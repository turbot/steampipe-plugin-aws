package aws

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go/middleware"
	"github.com/hashicorp/go-hclog"
	"github.com/turbot/steampipe-plugin-sdk/v6/connection"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/context_key"
)

// hydrateRateLimitWaiterKeyType is the type of the context key that
// plugin.WaitForHydrateRateLimit looks up. If the SDK renames it, the test
// below fails, and this name must be updated after checking that the SDK
// still reads the hydrate call's limiters from the request context.
const hydrateRateLimitWaiterKeyType = "plugin.hydrateRateLimitWaiterKey"

// waitRecordingContext records each lookup of the SDK's hydrate rate limit
// key, so a test can see that plugin.WaitForHydrateRateLimit ran with the
// request's context.
type waitRecordingContext struct {
	context.Context
	mu      sync.Mutex
	lookups int
}

func (c *waitRecordingContext) Value(key any) any {
	if fmt.Sprintf("%T", key) == hydrateRateLimitWaiterKeyType {
		c.mu.Lock()
		c.lookups++
		c.mu.Unlock()
	}
	return c.Context.Value(key)
}

func (c *waitRecordingContext) waits() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.lookups
}

var errProbeStop = errors.New("probe: stop before the request is sent")

// probeInitializeStep adds a middleware at the end of the Initialize step.
// It records the Initialize step's middleware IDs and stops the operation
// before it serializes the request, so no request is sent.
func probeInitializeStep(ids *[]string) func(*middleware.Stack) error {
	return func(stack *middleware.Stack) error {
		return stack.Initialize.Add(middleware.InitializeMiddlewareFunc("TestProbe",
			func(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler) (middleware.InitializeOutput, middleware.Metadata, error) {
				*ids = stack.Initialize.List()
				return middleware.InitializeOutput{}, middleware.Metadata{}, errProbeStop
			}), middleware.After)
	}
}

// Every service client the plugin builds must call
// plugin.WaitForHydrateRateLimit with the request's context before it sends
// a request. The plugin sets ClientWaitsForHydrateRateLimit, so the SDK does
// not wait before a hydrate call any more: a client that skips the wait sends
// its requests with no rate limit at all.
func TestAWSClientsWaitForHydrateRateLimit(t *testing.T) {
	isolateAWSEnv(t)
	conn := &plugin.Connection{Name: fmt.Sprintf("%s/%d", t.Name(), time.Now().UnixNano())}
	conn.SetConfig(awsConfig{
		AccessKey:        aws.String("AKIDEXAMPLE"),
		SecretKey:        aws.String("example-secret"),
		DefaultRegion:    aws.String("us-east-1"),
		Regions:          []string{"us-east-1"},
		S3ForcePathStyle: aws.Bool(true),
	})
	cache, err := connection.NewConnectionCache(conn.Name, 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	d := &plugin.QueryData{Connection: conn, ConnectionCache: cache}
	ctx := context.WithValue(context.Background(), context_key.Logger, hclog.NewNullLogger())

	// One client from each path a client config comes from: the default
	// region (getClientForDefaultRegion), a given region (getClient), fewer
	// retries (getClientWithMaxRetries), and a client with its own options
	// (S3 path-style addressing).
	cases := []struct {
		name string
		call func(ctx context.Context, probe func(*middleware.Stack) error) error
	}{
		{"STSClient", func(ctx context.Context, probe func(*middleware.Stack) error) error {
			svc, err := STSClient(ctx, d)
			if err != nil {
				return err
			}
			_, err = svc.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{}, func(o *sts.Options) { o.APIOptions = append(o.APIOptions, probe) })
			return err
		}},
		{"IAMClient", func(ctx context.Context, probe func(*middleware.Stack) error) error {
			svc, err := IAMClient(ctx, d)
			if err != nil {
				return err
			}
			_, err = svc.GetRole(ctx, &iam.GetRoleInput{RoleName: aws.String("r")}, func(o *iam.Options) { o.APIOptions = append(o.APIOptions, probe) })
			return err
		}},
		{"EC2ClientForRegion", func(ctx context.Context, probe func(*middleware.Stack) error) error {
			svc, err := EC2ClientForRegion(ctx, d, "us-east-1")
			if err != nil {
				return err
			}
			_, err = svc.DescribeSnapshotAttribute(ctx, &ec2.DescribeSnapshotAttributeInput{SnapshotId: aws.String("snap-1"), Attribute: ec2types.SnapshotAttributeNameCreateVolumePermission}, func(o *ec2.Options) { o.APIOptions = append(o.APIOptions, probe) })
			return err
		}},
		{"EC2LowRetryClientForRegion", func(ctx context.Context, probe func(*middleware.Stack) error) error {
			svc, err := EC2LowRetryClientForRegion(ctx, d, "us-east-1")
			if err != nil {
				return err
			}
			_, err = svc.DescribeRegions(ctx, &ec2.DescribeRegionsInput{}, func(o *ec2.Options) { o.APIOptions = append(o.APIOptions, probe) })
			return err
		}},
		{"S3Client", func(ctx context.Context, probe func(*middleware.Stack) error) error {
			svc, err := S3Client(ctx, d, "us-east-1")
			if err != nil {
				return err
			}
			_, err = svc.GetBucketPolicy(ctx, &s3.GetBucketPolicyInput{Bucket: aws.String("b")}, func(o *s3.Options) { o.APIOptions = append(o.APIOptions, probe) })
			return err
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var ids []string
			rctx := &waitRecordingContext{Context: ctx}
			if err := tc.call(rctx, probeInitializeStep(&ids)); !errors.Is(err, errProbeStop) {
				t.Fatalf("operation error = %v, want the probe to stop it", err)
			}
			if got := rctx.waits(); got != 1 {
				t.Errorf("plugin.WaitForHydrateRateLimit ran %d times before the request, want 1", got)
			}
			if !slices.Contains(ids, hydrateRateLimitMiddlewareID) {
				t.Errorf("Initialize step = %v, want %s in it", ids, hydrateRateLimitMiddlewareID)
			}
		})
	}
}

// The SDK only skips its own wait before a hydrate call when the plugin says
// that its clients wait. Without it, a hydrate call that sends no request
// still takes a token, e.g., create_volume_permissions for public snapshots.
func TestPluginClientWaitsForHydrateRateLimit(t *testing.T) {
	if !Plugin(context.Background()).ClientWaitsForHydrateRateLimit {
		t.Fatal("Plugin.ClientWaitsForHydrateRateLimit = false, want true")
	}
}

// The test above covers the client paths that exist today. This guards the
// next one: every client config must come from getBaseClientForAccountUncached,
// which adds the rate limit middleware, so no other code may load or build an
// aws.Config, or build a service client from bare options.
func TestOnlyBaseClientBuildsAWSConfig(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	const allowed = "getBaseClientForAccountUncached"
	fset := token.NewFileSet()
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, file, nil, parser.SkipObjectResolution)
		if err != nil {
			t.Fatal(err)
		}
		imports := map[string]string{}
		for _, imp := range f.Imports {
			path, _ := strconv.Unquote(imp.Path.Value)
			name := filepath.Base(path)
			if imp.Name != nil {
				name = imp.Name.Name
			}
			imports[name] = path
		}
		// selector reports whether expr is pkg.sel, with pkg imported from path.
		selector := func(expr ast.Expr, sel string, match func(path string) bool) bool {
			s, ok := expr.(*ast.SelectorExpr)
			if !ok || s.Sel.Name != sel {
				return false
			}
			x, ok := s.X.(*ast.Ident)
			return ok && match(imports[x.Name])
		}
		is := func(path string) func(string) bool { return func(p string) bool { return p == path } }
		isService := func(p string) bool { return strings.HasPrefix(p, "github.com/aws/aws-sdk-go-v2/service/") }

		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil || fn.Name.Name == allowed {
				continue
			}
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				var what string
				switch n := n.(type) {
				case *ast.CallExpr:
					switch {
					case selector(n.Fun, "LoadDefaultConfig", is("github.com/aws/aws-sdk-go-v2/config")):
						what = "config.LoadDefaultConfig"
					case selector(n.Fun, "NewConfig", is("github.com/aws/aws-sdk-go-v2/aws")):
						what = "aws.NewConfig"
					}
				case *ast.CompositeLit:
					switch {
					case selector(n.Type, "Config", is("github.com/aws/aws-sdk-go-v2/aws")):
						what = "an aws.Config literal"
					case selector(n.Type, "Options", isService):
						what = "a service client Options literal"
					}
				}
				if what != "" {
					t.Errorf("%s: %s builds %s. Build clients from the config of getClient (or a helper that uses it), so they wait for the hydrate call's rate limiters.",
						fset.Position(n.Pos()), fn.Name.Name, what)
				}
				return true
			})
		}
	}
}
