package aws

import (
	"context"
	"testing"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// TestConnectionConfigCredentialsProvider_PicksUpInPlaceConfigMutation
// verifies that the provider returns the CURRENT contents of
// Connection.Config on every Retrieve call, not values captured at
// construction time. This is the core property that lets in-flight goroutines
// pick up rotated credentials.
func TestConnectionConfigCredentialsProvider_PicksUpInPlaceConfigMutation(t *testing.T) {
	ak1, sk1, st1 := "AKIA1ORIGINAL", "secret1original", "token1original"
	ak2, sk2, st2 := "AKIA2ROTATED", "secret2rotated", "token2rotated"

	conn := &plugin.Connection{
		Name: "test",
		Config: awsConfig{
			AccessKey:    &ak1,
			SecretKey:    &sk1,
			SessionToken: &st1,
		},
	}
	provider := &connectionConfigCredentialsProvider{connection: conn}

	got1, err := provider.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("first Retrieve failed: %v", err)
	}
	if got1.AccessKeyID != ak1 || got1.SecretAccessKey != sk1 || got1.SessionToken != st1 {
		t.Errorf("first Retrieve returned wrong creds: got AK=%s SK=%s ST=%s, want AK=%s SK=%s ST=%s",
			got1.AccessKeyID, got1.SecretAccessKey, got1.SessionToken, ak1, sk1, st1)
	}

	// Mutate the connection's Config in place, mimicking what the plugin SDK
	// does at steampipe-plugin-sdk/plugin/plugin_connection_config.go:154
	// (d.Connection.Config = configStruct) when UpdateConnectionConfigs
	// delivers rotated credentials.
	conn.Config = awsConfig{
		AccessKey:    &ak2,
		SecretKey:    &sk2,
		SessionToken: &st2,
	}

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
	conn := &plugin.Connection{
		Name: "test",
		Config: awsConfig{
			AccessKey: &ak,
		},
	}
	provider := &connectionConfigCredentialsProvider{connection: conn}
	if _, err := provider.Retrieve(context.Background()); err == nil {
		t.Errorf("Retrieve should error when secret_key is missing, got nil")
	}
}

func TestConnectionConfigCredentialsProvider_ErrorsWhenAccessKeyMissing(t *testing.T) {
	sk := "secret"
	conn := &plugin.Connection{
		Name: "test",
		Config: awsConfig{
			SecretKey: &sk,
		},
	}
	provider := &connectionConfigCredentialsProvider{connection: conn}
	if _, err := provider.Retrieve(context.Background()); err == nil {
		t.Errorf("Retrieve should error when access_key is missing, got nil")
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
	conn := &plugin.Connection{
		Name: "test",
		Config: awsConfig{
			AccessKey: &ak,
			SecretKey: &sk,
		},
	}
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

func TestConnectionConfigCredentialsProvider_SetsCanExpireForCacheRefresh(t *testing.T) {
	// The SDK's CredentialsCache will only call Retrieve again before Expires
	// if CanExpire is true. Without this, in-flight goroutines wouldn't pick
	// up rotated creds even with our dynamic provider.
	ak, sk, st := "AKIA1", "secret1", "token1"
	conn := &plugin.Connection{
		Name: "test",
		Config: awsConfig{
			AccessKey:    &ak,
			SecretKey:    &sk,
			SessionToken: &st,
		},
	}
	provider := &connectionConfigCredentialsProvider{connection: conn}
	got, err := provider.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("Retrieve failed: %v", err)
	}
	if !got.CanExpire {
		t.Errorf("Credentials.CanExpire should be true so the SDK CredentialsCache calls Retrieve again before TTL")
	}
	if got.Expires.IsZero() {
		t.Errorf("Credentials.Expires should be set to a non-zero hint")
	}
}
