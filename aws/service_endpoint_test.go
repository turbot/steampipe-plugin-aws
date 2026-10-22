package aws

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/hashicorp/go-hclog"
	"github.com/turbot/steampipe-plugin-sdk/v6/connection"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/context_key"
)

// endpointStub stands in for an AWS endpoint. It answers sts:GetCallerIdentity
// with its own account ID, answers sts:AssumeRole with the key ASSUMEDKEY, and
// records the access key ID of each request.
type endpointStub struct {
	*httptest.Server
	account string
	mu      sync.Mutex
	keys    []string
}

func newEndpointStub(t *testing.T, account string) *endpointStub {
	t.Helper()
	s := &endpointStub{account: account}
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, cred, _ := strings.Cut(r.Header.Get("Authorization"), "Credential=")
		key, _, _ := strings.Cut(cred, "/")
		s.mu.Lock()
		s.keys = append(s.keys, key)
		s.mu.Unlock()
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/xml")
		if r.PostForm.Get("Action") == "AssumeRole" {
			_, _ = fmt.Fprintf(w, `<AssumeRoleResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/">
  <AssumeRoleResult>
    <Credentials><AccessKeyId>ASSUMEDKEY</AccessKeyId><SecretAccessKey>assumed-secret</SecretAccessKey><SessionToken>assumed-token</SessionToken><Expiration>2099-01-01T00:00:00Z</Expiration></Credentials>
    <AssumedRoleUser><Arn>arn:aws:sts::%[1]s:assumed-role/reader/session</Arn><AssumedRoleId>AROAEXAMPLE:session</AssumedRoleId></AssumedRoleUser>
  </AssumeRoleResult>
  <ResponseMetadata><RequestId>00000000-0000-0000-0000-000000000000</RequestId></ResponseMetadata>
</AssumeRoleResponse>`, s.account)
			return
		}
		_, _ = fmt.Fprintf(w, `<GetCallerIdentityResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/">
  <GetCallerIdentityResult><Arn>arn:aws:iam::%[1]s:root</Arn><UserId>%[1]s</UserId><Account>%[1]s</Account></GetCallerIdentityResult>
  <ResponseMetadata><RequestId>00000000-0000-0000-0000-000000000000</RequestId></ResponseMetadata>
</GetCallerIdentityResponse>`, s.account)
	}))
	t.Cleanup(s.Close)
	return s
}

func (s *endpointStub) seenKeys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.keys)
}

// isolateAWSEnv removes AWS settings of the test process and of the machine,
// so each test controls every source of endpoint and credentials.
func isolateAWSEnv(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("AWS_CONFIG_FILE", filepath.Join(dir, "no-config"))
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", filepath.Join(dir, "no-credentials"))
	t.Setenv("AWS_EC2_METADATA_DISABLED", "true")
	for _, k := range []string{
		"AWS_ACCESS_KEY_ID", "AWS_ACCESS_KEY", "AWS_SECRET_ACCESS_KEY", "AWS_SECRET_KEY", "AWS_SESSION_TOKEN",
		"AWS_PROFILE", "AWS_DEFAULT_PROFILE", "AWS_REGION", "AWS_DEFAULT_REGION",
		"AWS_ROLE_ARN", "AWS_WEB_IDENTITY_TOKEN_FILE",
		"AWS_CONTAINER_CREDENTIALS_FULL_URI", "AWS_CONTAINER_CREDENTIALS_RELATIVE_URI",
		"AWS_ENDPOINT_URL", "AWS_ENDPOINT_URL_STS", "AWS_IGNORE_CONFIGURED_ENDPOINT_URLS",
	} {
		t.Setenv(k, "")
		if err := os.Unsetenv(k); err != nil {
			t.Fatal(err)
		}
	}
}

// callerAccount calls sts:GetCallerIdentity through the plugin's own client
// path for a connection with the given config. The connection name is unique
// per call, because the plugin caches clients per connection name.
func callerAccount(t *testing.T, name string, cfg awsConfig) (string, error) {
	t.Helper()
	cfg.DefaultRegion = aws.String("us-east-1")
	cfg.Regions = []string{"us-east-1"}
	cfg.MaxErrorRetryAttempts = aws.Int(1)

	conn := &plugin.Connection{Name: t.Name() + "/" + name}
	conn.SetConfig(cfg)
	cache, err := connection.NewConnectionCache(conn.Name, 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	d := &plugin.QueryData{Connection: conn, ConnectionCache: cache}
	ctx := context.WithValue(context.Background(), context_key.Logger, hclog.NewNullLogger())

	svc, err := STSClient(ctx, d)
	if err != nil {
		return "", err
	}
	out, err := svc.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}
	return aws.ToString(out.Account), nil
}

func wantKeys(t *testing.T, s *endpointStub, want ...string) {
	t.Helper()
	got := s.seenKeys()
	if want == nil {
		want = []string{}
	}
	if got == nil {
		got = []string{}
	}
	if !slices.Equal(got, want) {
		t.Errorf("endpoint %s saw access keys %q, want %q", s.URL, got, want)
	}
}

// Two connections in one process, each with its own endpoint_url and keys.
// Each must reach its own endpoint and sign with its own key. The process has
// keys and AWS_ENDPOINT_URL of its own, which must not be used.
func TestEndpointUrlSignsWithConnectionKeys(t *testing.T) {
	for _, processKeys := range []bool{false, true} {
		t.Run(fmt.Sprintf("process_keys=%t", processKeys), func(t *testing.T) {
			isolateAWSEnv(t)
			one := newEndpointStub(t, "111111111111")
			two := newEndpointStub(t, "222222222222")
			trap := newEndpointStub(t, "999999999999")
			t.Setenv("AWS_ENDPOINT_URL", trap.URL)
			if processKeys {
				t.Setenv("AWS_ACCESS_KEY_ID", "PROCESSKEY")
				t.Setenv("AWS_SECRET_ACCESS_KEY", "process-secret")
			}

			for _, c := range []struct {
				stub *endpointStub
				key  string
			}{{one, "CONNECTIONKEYONE"}, {two, "CONNECTIONKEYTWO"}} {
				account, err := callerAccount(t, c.key, awsConfig{
					EndpointUrl: aws.String(c.stub.URL),
					AccessKey:   aws.String(c.key),
					SecretKey:   aws.String("secret-" + c.key),
				})
				if err != nil {
					t.Fatalf("connection with key %s: %v", c.key, err)
				}
				if account != c.stub.account {
					t.Errorf("connection with key %s got account %s, want %s", c.key, account, c.stub.account)
				}
			}

			wantKeys(t, one, "CONNECTIONKEYONE")
			wantKeys(t, two, "CONNECTIONKEYTWO")
			wantKeys(t, trap)
		})
	}
}

// A connection with endpoint_url and a profile must use the profile's
// credentials.
func TestEndpointUrlUsesConnectionProfile(t *testing.T) {
	isolateAWSEnv(t)
	stub := newEndpointStub(t, "333333333333")
	configFile := filepath.Join(t.TempDir(), "config")
	err := os.WriteFile(configFile, []byte("[profile emulator]\naws_access_key_id = PROFILEKEY\naws_secret_access_key = profile-secret\n"), 0o600)
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("AWS_CONFIG_FILE", configFile)

	_, err = callerAccount(t, "profile", awsConfig{
		EndpointUrl: aws.String(stub.URL),
		Profile:     aws.String("emulator"),
	})
	if err != nil {
		t.Fatal(err)
	}
	wantKeys(t, stub, "PROFILEKEY")
}

// A connection with endpoint_url and a role profile must send the SDK's own
// sts:AssumeRole to the connection's endpoint too, and then sign with the
// role's keys.
func TestEndpointUrlAssumesRoleAtConnectionEndpoint(t *testing.T) {
	isolateAWSEnv(t)
	stub := newEndpointStub(t, "666666666666")
	configFile := filepath.Join(t.TempDir(), "config")
	err := os.WriteFile(configFile, []byte(`[profile source]
aws_access_key_id = SOURCEKEY
aws_secret_access_key = source-secret

[profile reader]
role_arn = arn:aws:iam::666666666666:role/reader
source_profile = source
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("AWS_CONFIG_FILE", configFile)

	_, err = callerAccount(t, "role", awsConfig{
		EndpointUrl: aws.String(stub.URL),
		Profile:     aws.String("reader"),
	})
	if err != nil {
		t.Fatal(err)
	}
	wantKeys(t, stub, "SOURCEKEY", "ASSUMEDKEY")
}

// endpoint_url is used for all AWS API requests of its connection, so a
// service-specific endpoint of the process must not override it.
func TestEndpointUrlBeatsServiceEndpointEnv(t *testing.T) {
	isolateAWSEnv(t)
	stub := newEndpointStub(t, "444444444444")
	trap := newEndpointStub(t, "999999999999")
	t.Setenv("AWS_ENDPOINT_URL_STS", trap.URL)

	_, err := callerAccount(t, "service-env", awsConfig{
		EndpointUrl: aws.String(stub.URL),
		AccessKey:   aws.String("CONNECTIONKEY"),
		SecretKey:   aws.String("connection-secret"),
	})
	if err != nil {
		t.Fatal(err)
	}
	wantKeys(t, stub, "CONNECTIONKEY")
	wantKeys(t, trap)
}

// With no endpoint_url, the AWS_ENDPOINT_URL environment variable is used, as
// the plugin docs say.
func TestEndpointUrlFromEnvWhenNotConfigured(t *testing.T) {
	isolateAWSEnv(t)
	stub := newEndpointStub(t, "555555555555")
	t.Setenv("AWS_ENDPOINT_URL", stub.URL)

	_, err := callerAccount(t, "env", awsConfig{
		AccessKey: aws.String("CONNECTIONKEY"),
		SecretKey: aws.String("connection-secret"),
	})
	if err != nil {
		t.Fatal(err)
	}
	wantKeys(t, stub, "CONNECTIONKEY")
}
