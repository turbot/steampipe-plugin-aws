package aws

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/hashicorp/go-hclog"
	"github.com/turbot/steampipe-plugin-sdk/v6/connection"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/context_key"
)

const maxIdleConnsEnv = "STEAMPIPE_AWS_HTTP_TRANSPORT_MAX_IDLE_CONNS"

// newCountingSTSStub starts an endpoint that answers sts:GetCallerIdentity
// over HTTP/1.1, like the AWS API endpoints, so each request in flight needs
// its own connection. (Plain HTTP pools connections the same way as HTTPS.)
// It holds each request for a short time, so requests sent together overlap.
// It counts the TCP connections it accepts.
func newCountingSTSStub(t *testing.T) (*httptest.Server, *atomic.Int64) {
	t.Helper()
	var accepted atomic.Int64
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(20 * time.Millisecond)
		w.Header().Set("Content-Type", "text/xml")
		_, _ = io.WriteString(w, `<GetCallerIdentityResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/">
  <GetCallerIdentityResult><Arn>arn:aws:iam::111111111111:root</Arn><UserId>111111111111</UserId><Account>111111111111</Account></GetCallerIdentityResult>
  <ResponseMetadata><RequestId>00000000-0000-0000-0000-000000000000</RequestId></ResponseMetadata>
</GetCallerIdentityResponse>`)
	}))
	srv.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		if state == http.StateNew {
			accepted.Add(1)
		}
	}
	srv.Start()
	t.Cleanup(srv.Close)
	return srv, &accepted
}

// closeCountingConn keeps a count of the connections that are still open.
type closeCountingConn struct {
	net.Conn
	once sync.Once
	open *atomic.Int64
}

func (c *closeCountingConn) Close() error {
	c.once.Do(func() { c.open.Add(-1) })
	return c.Conn.Close()
}

// A hydrate function makes a new service client for each call, and a table
// like aws_sns_topic sends a burst of these calls to one host for each page
// of rows. All the clients share sharedHTTPClient, so each burst must reuse
// the connections of the bursts before it. With the AWS SDK default of 10
// idle connections per host, every connection above 10 was closed when a
// burst ended, so each later burst of 200 opened about 190 new connections
// (#2816).
func TestPluginClientsReuseConnectionsAcrossBursts(t *testing.T) {
	if os.Getenv(maxIdleConnsEnv) != "" {
		t.Skipf("%s is set, and this test checks the default", maxIdleConnsEnv)
	}
	isolateAWSEnv(t)
	// With any other defaults mode, the AWS SDK gives each service client its
	// own copy of the HTTP client, and so its own empty pool (#2816).
	t.Setenv("AWS_DEFAULTS_MODE", "")
	srv, accepted := newCountingSTSStub(t)

	// The plugin caches clients per connection name, so make the name unique
	// for each run of the test.
	conn := &plugin.Connection{Name: fmt.Sprintf("%s/%d", t.Name(), time.Now().UnixNano())}
	conn.SetConfig(awsConfig{
		AccessKey:             aws.String("AKIDEXAMPLE"),
		SecretKey:             aws.String("example-secret"),
		EndpointUrl:           aws.String(srv.URL),
		DefaultRegion:         aws.String("us-east-1"),
		Regions:               []string{"us-east-1"},
		MaxErrorRetryAttempts: aws.Int(1),
	})
	cache, err := connection.NewConnectionCache(conn.Name, 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	d := &plugin.QueryData{Connection: conn, ConnectionCache: cache}
	ctx := context.WithValue(context.Background(), context_key.Logger, hclog.NewNullLogger())

	const burst, bursts = 200, 5
	for range bursts {
		var wg sync.WaitGroup
		for range burst {
			wg.Add(1)
			go func() {
				defer wg.Done()
				// A new client for each call, as in a hydrate function.
				svc, err := STSClient(ctx, d)
				if err != nil {
					t.Error(err)
					return
				}
				if _, err := svc.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{}); err != nil {
					t.Error(err)
				}
			}()
		}
		wg.Wait()
		// The transport puts a connection back in its idle pool just after
		// the response is read, so give it time before the next burst.
		time.Sleep(100 * time.Millisecond)
	}

	// The first burst needs 200 connections. The later bursts must reuse them.
	if got := accepted.Load(); got > 2*burst {
		t.Errorf("%d bursts of %d requests to one host opened %d connections, want at most %d", bursts, burst, got, 2*burst)
	}
}

// The idle pool must have a total cap. aws_s3_bucket sends requests to two
// new hosts for each bucket and never uses them again. Without a total cap,
// the connections of every bucket would stay open until the idle timeout.
func TestHTTPClientCapsIdleConnectionsAcrossHosts(t *testing.T) {
	t.Setenv(maxIdleConnsEnv, "20")
	srv, _ := newCountingSTSStub(t)
	var open atomic.Int64
	client := initializeHTTPClient().(*awshttp.BuildableClient).WithTransportOptions(func(tr *http.Transport) {
		// Send every host to the test server, and count the open connections.
		tr.DialContext = func(ctx context.Context, network, _ string) (net.Conn, error) {
			conn, err := (&net.Dialer{}).DialContext(ctx, network, srv.Listener.Addr().String())
			if err != nil {
				return nil, err
			}
			open.Add(1)
			return &closeCountingConn{Conn: conn, open: &open}, nil
		}
	})

	_, port, err := net.SplitHostPort(srv.Listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for bucket := range 60 {
		for range 3 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://bucket-%d.s3.example.com:%s/", bucket, port), nil)
				if err != nil {
					t.Error(err)
					return
				}
				resp, err := client.Do(req)
				if err != nil {
					t.Error(err)
					return
				}
				// Read and close the body, as the AWS SDK does, so the
				// connection can go back to the idle pool.
				_, _ = io.Copy(io.Discard, resp.Body)
				_ = resp.Body.Close()
			}()
		}
	}
	wg.Wait()

	// The transport closes an idle connection over the cap when another one
	// comes back to the pool, just after each request ends.
	for deadline := time.Now().Add(2 * time.Second); open.Load() > 20 && time.Now().Before(deadline); {
		time.Sleep(10 * time.Millisecond)
	}
	if got := open.Load(); got > 20 {
		t.Errorf("after requests to 60 hosts, %d idle connections are open, want at most 20 (%s=20)", got, maxIdleConnsEnv)
	}
}

func TestHTTPClientIdleConnectionSettings(t *testing.T) {
	cases := []struct {
		env                   string
		wantTotal, wantPerHst int
	}{
		{env: "", wantTotal: 1000, wantPerHst: 1000},
		{env: "50", wantTotal: 50, wantPerHst: 50},
		// 0 keeps the AWS SDK defaults.
		{env: "0", wantTotal: awshttp.DefaultHTTPTransportMaxIdleConns, wantPerHst: awshttp.DefaultHTTPTransportMaxIdleConnsPerHost},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%s=%q", maxIdleConnsEnv, c.env), func(t *testing.T) {
			t.Setenv(maxIdleConnsEnv, c.env)
			t.Setenv("STEAMPIPE_AWS_HTTP_TRANSPORT_MAX_CONNS_PER_HOST", "")
			tr := initializeHTTPClient().(*awshttp.BuildableClient).GetTransport()
			if tr.MaxIdleConns != c.wantTotal || tr.MaxIdleConnsPerHost != c.wantPerHst {
				t.Errorf("MaxIdleConns=%d MaxIdleConnsPerHost=%d, want %d and %d", tr.MaxIdleConns, tr.MaxIdleConnsPerHost, c.wantTotal, c.wantPerHst)
			}
			if tr.MaxConnsPerHost != 5000 {
				t.Errorf("MaxConnsPerHost=%d, want 5000 (unchanged)", tr.MaxConnsPerHost)
			}
		})
	}
}
