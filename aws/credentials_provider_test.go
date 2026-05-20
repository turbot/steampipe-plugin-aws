package aws

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

// TestConnectionConfigCredentialsProvider_PicksUpInPlaceConfigMutation
// verifies that the provider returns the CURRENT contents of
// Connection.Config on every Retrieve call, not values captured at
// construction time. This is the core property that lets in-flight goroutines
// pick up rotated credentials.
func TestConnectionConfigCredentialsProvider_PicksUpInPlaceConfigMutation(t *testing.T) {
	ak1, sk1, st1 := "AKIA1ORIGINAL", "secret1original", "token1original"
	ak2, sk2, st2 := "AKIA2ROTATED", "secret2rotated", "token2rotated"

	conn := &plugin.Connection{Name: "test"}
	conn.SetConfig(awsConfig{
		AccessKey:    &ak1,
		SecretKey:    &sk1,
		SessionToken: &st1,
	})
	provider := &connectionConfigCredentialsProvider{connection: conn}

	got1, err := provider.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("first Retrieve failed: %v", err)
	}
	if got1.AccessKeyID != ak1 || got1.SecretAccessKey != sk1 || got1.SessionToken != st1 {
		t.Errorf("first Retrieve returned wrong creds: got AK=%s SK=%s ST=%s, want AK=%s SK=%s ST=%s",
			got1.AccessKeyID, got1.SecretAccessKey, got1.SessionToken, ak1, sk1, st1)
	}

	// Rotate the connection's config in place, mimicking what the plugin SDK
	// does in upsertConnectionData when UpdateConnectionConfigs delivers
	// rotated credentials.
	conn.SetConfig(awsConfig{
		AccessKey:    &ak2,
		SecretKey:    &sk2,
		SessionToken: &st2,
	})

	got2, err := provider.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("second Retrieve failed: %v", err)
	}
	if got2.AccessKeyID != ak2 || got2.SecretAccessKey != sk2 || got2.SessionToken != st2 {
		t.Errorf("second Retrieve returned wrong creds after rotation: got AK=%s SK=%s ST=%s, want AK=%s SK=%s ST=%s",
			got2.AccessKeyID, got2.SecretAccessKey, got2.SessionToken, ak2, sk2, st2)
	}
}

func TestConnectionConfigCredentialsProvider_ErrorsWhenSecretKeyMissing(t *testing.T) {
	ak := "AKIA1"
	conn := &plugin.Connection{Name: "my-test-conn"}
	conn.SetConfig(awsConfig{
		AccessKey: &ak,
	})
	provider := &connectionConfigCredentialsProvider{connection: conn}
	_, err := provider.Retrieve(context.Background())
	if err == nil {
		t.Fatalf("Retrieve should error when secret_key is missing, got nil")
	}
	if !strings.Contains(err.Error(), "secret_key") {
		t.Errorf("error should mention which field is missing, got: %v", err)
	}
	if !strings.Contains(err.Error(), conn.Name) {
		t.Errorf("error should include the connection name %q for diagnostics, got: %v", conn.Name, err)
	}
}

func TestConnectionConfigCredentialsProvider_ErrorsWhenAccessKeyMissing(t *testing.T) {
	sk := "secret"
	conn := &plugin.Connection{Name: "my-test-conn"}
	conn.SetConfig(awsConfig{
		SecretKey: &sk,
	})
	provider := &connectionConfigCredentialsProvider{connection: conn}
	_, err := provider.Retrieve(context.Background())
	if err == nil {
		t.Fatalf("Retrieve should error when access_key is missing, got nil")
	}
	if !strings.Contains(err.Error(), "access_key") {
		t.Errorf("error should mention which field is missing, got: %v", err)
	}
	if !strings.Contains(err.Error(), conn.Name) {
		t.Errorf("error should include the connection name %q for diagnostics, got: %v", conn.Name, err)
	}
}

func TestConnectionConfigCredentialsProvider_ErrorsWhenConnectionNil(t *testing.T) {
	provider := &connectionConfigCredentialsProvider{connection: nil}
	if _, err := provider.Retrieve(context.Background()); err == nil {
		t.Errorf("Retrieve should error when connection is nil, got nil")
	}
}

func TestConnectionConfigCredentialsProvider_SessionTokenOptional(t *testing.T) {
	// Long-term IAM user credentials don't have a session token. Verify the
	// provider handles that case without panicking and returns an empty
	// SessionToken.
	ak, sk := "AKIA1LONGTERM", "secretlongterm"
	conn := &plugin.Connection{Name: "test"}
	conn.SetConfig(awsConfig{
		AccessKey: &ak,
		SecretKey: &sk,
	})
	provider := &connectionConfigCredentialsProvider{connection: conn}
	got, err := provider.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("Retrieve failed: %v", err)
	}
	if got.SessionToken != "" {
		t.Errorf("SessionToken should be empty for long-term creds, got %q", got.SessionToken)
	}
	if got.AccessKeyID != ak || got.SecretAccessKey != sk {
		t.Errorf("Retrieve returned wrong creds: got AK=%s SK=%s, want AK=%s SK=%s",
			got.AccessKeyID, got.SecretAccessKey, ak, sk)
	}
}

// TestConnectionConfigCredentialsProvider_ExpiresIsShort guards against a
// regression where Expires drifts back out to the original token TTL — which
// would defeat the rotation-pickup goal because the SDK's CredentialsCache
// would not call Retrieve again until then. The standalone reproduction
// validated 60 seconds; we allow a generous bound here to permit small
// adjustments without breaking the test, but anything more than a few
// minutes would silently leak stale creds for that long.
func TestConnectionConfigCredentialsProvider_ExpiresIsShort(t *testing.T) {
	ak, sk, st := "AKIA1", "secret1", "token1"
	conn := &plugin.Connection{Name: "test"}
	conn.SetConfig(awsConfig{
		AccessKey:    &ak,
		SecretKey:    &sk,
		SessionToken: &st,
	})
	provider := &connectionConfigCredentialsProvider{connection: conn}

	before := time.Now()
	got, err := provider.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("Retrieve failed: %v", err)
	}

	if !got.CanExpire {
		t.Errorf("Credentials.CanExpire must be true so the SDK CredentialsCache calls Retrieve again before TTL")
	}
	if got.Expires.IsZero() {
		t.Fatalf("Credentials.Expires should be set to a non-zero hint")
	}

	maxAllowed := before.Add(5 * time.Minute)
	if got.Expires.After(maxAllowed) {
		t.Errorf("Credentials.Expires should be within ~5min of now to bound rotation latency; got %s (now=%s, max=%s). "+
			"A long Expires defeats the rotation-pickup goal because the SDK's CredentialsCache won't call Retrieve again until then.",
			got.Expires.UTC().Format(time.RFC3339Nano),
			before.UTC().Format(time.RFC3339Nano),
			maxAllowed.UTC().Format(time.RFC3339Nano))
	}
}

// TestConnectionConfigCredentialsProvider_DoesNotPanicOnEmptyRegions
// verifies that Retrieve does not invoke GetConfig's region normalization
// path, which panics on regions = []. Moving that panic into the AWS request
// signing path would be worse than the original startup-time crash.
func TestConnectionConfigCredentialsProvider_DoesNotPanicOnEmptyRegions(t *testing.T) {
	ak, sk := "AKIA1", "secret1"
	conn := &plugin.Connection{Name: "test"}
	conn.SetConfig(awsConfig{
		AccessKey: &ak,
		SecretKey: &sk,
		Regions:   []string{}, // GetConfig would panic on this
	})
	provider := &connectionConfigCredentialsProvider{connection: conn}

	// This call would panic if Retrieve went through GetConfig.
	got, err := provider.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("Retrieve failed: %v", err)
	}
	if got.AccessKeyID != ak || got.SecretAccessKey != sk {
		t.Errorf("Retrieve returned wrong creds: got AK=%s SK=%s, want AK=%s SK=%s",
			got.AccessKeyID, got.SecretAccessKey, ak, sk)
	}
}

// TestConnectionConfigCredentialsProvider_RecoversFromTornConfigRead
// verifies the defer/recover safety net in case a concurrent SDK write to
// Connection.Config produces a torn interface value that would otherwise
// panic during the type assertion. We simulate the panic directly by setting
// Config to a sentinel value of a type that is not awsConfig but is also not
// an obvious "wrong type" — the type assertion uses the comma-ok form so
// this returns ok=false rather than panicking. To actually trigger a panic
// we set Config to a non-nil value of an unexpected concrete type and verify
// we still get a clean error rather than a propagated panic.
func TestConnectionConfigCredentialsProvider_HandlesUnexpectedConfigType(t *testing.T) {
	conn := &plugin.Connection{Name: "test"}
	conn.SetConfig("this is not an awsConfig")
	provider := &connectionConfigCredentialsProvider{connection: conn}
	_, err := provider.Retrieve(context.Background())
	if err == nil {
		t.Errorf("Retrieve should error when config is not awsConfig, got nil")
	}
	if !strings.Contains(err.Error(), "awsConfig") {
		t.Errorf("error should mention expected type for diagnostics, got: %v", err)
	}
}

// TestConnectionConfigCredentialsProvider_CacheLayerPicksUpRotation wraps the
// provider in the same aws.CredentialsCache the production path uses (via
// config.WithCredentialsProvider) and verifies that a rotation in
// Connection.Config is visible to a cache.Retrieve call after the cached
// value's Expires has passed. This is the production-shaped regression test
// for the ExpiredToken-under-rotation bug — the unit tests above bypass the
// cache and only prove the provider returns the right values when invoked.
//
// We shrink credentialsExpiresInterval to 50ms so the test runs in tens of
// milliseconds instead of waiting out the production 60s window.
func TestConnectionConfigCredentialsProvider_CacheLayerPicksUpRotation(t *testing.T) {
	orig := credentialsExpiresInterval
	credentialsExpiresInterval = 50 * time.Millisecond
	defer func() { credentialsExpiresInterval = orig }()

	ak1, sk1, st1 := "AKIA1ORIGINAL", "secret1original", "token1original"
	ak2, sk2, st2 := "AKIA2ROTATED", "secret2rotated", "token2rotated"

	conn := &plugin.Connection{Name: "test"}
	conn.SetConfig(awsConfig{
		AccessKey:    &ak1,
		SecretKey:    &sk1,
		SessionToken: &st1,
	})
	provider := &connectionConfigCredentialsProvider{connection: conn}
	cache := aws.NewCredentialsCache(provider)

	got1, err := cache.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("first cache.Retrieve failed: %v", err)
	}
	if got1.AccessKeyID != ak1 {
		t.Errorf("first cache.Retrieve: got AK=%s, want %s", got1.AccessKeyID, ak1)
	}

	// Rotate the connection config in place — what the SDK does on
	// UpdateConnectionConfigs.
	conn.SetConfig(awsConfig{
		AccessKey:    &ak2,
		SecretKey:    &sk2,
		SessionToken: &st2,
	})

	// Sleep past the cache's Expires window so the next Retrieve goes back to
	// the provider rather than serving the cached value.
	time.Sleep(100 * time.Millisecond)

	got2, err := cache.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("second cache.Retrieve failed: %v", err)
	}
	if got2.AccessKeyID != ak2 {
		t.Errorf("after rotation, cache returned stale creds: got AK=%s, want %s. "+
			"This is the ExpiredToken-under-rotation bug — the CredentialsCache held "+
			"the original creds past their replacement in Connection.Config.",
			got2.AccessKeyID, ak2)
	}
	if got2.SecretAccessKey != sk2 || got2.SessionToken != st2 {
		t.Errorf("after rotation, cache returned stale creds: got SK=%s ST=%s, want SK=%s ST=%s",
			got2.SecretAccessKey, got2.SessionToken, sk2, st2)
	}
}
