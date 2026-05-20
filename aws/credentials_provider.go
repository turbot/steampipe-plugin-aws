package aws

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

// connectionConfigCredentialsProvider is an aws.CredentialsProvider that reads
// access_key / secret_key / session_token from the steampipe connection config
// on every Retrieve() call (subject to the SDK's CredentialsCache, which we
// hint at via a short Expires below).
//
// Background
//
// The steampipe plugin SDK mutates the connection config in place when a new
// connection config arrives via UpdateConnectionConfigs (see
// steampipe-plugin-sdk/plugin/plugin_connection_config.go upsertConnectionData,
// which calls d.Connection.SetConfig(configStruct) under a write lock). The
// SDK comment at that site explicitly acknowledges that a query may already be
// executing with this Connection object, and that the AWS plugin in particular
// "may refresh the Client using the previous credentials" — which is the bug
// this provider fixes.
//
// A credentials.NewStaticCredentialsProvider built once at aws.Config
// construction time captures the original token values. A goroutine holding
// that aws.Config keeps signing requests with the original token regardless of
// rotation. When the original token expires at AWS, every subsequent request
// from that goroutine fails with ExpiredToken — even if a fresh valid token
// has been delivered to Connection.Config by the SDK.
//
// By re-reading the connection config on every Retrieve(), in-flight goroutines
// holding the same aws.Config pick up rotated credentials on the next signing
// operation (modulo the CredentialsCache TTL we set below).
type connectionConfigCredentialsProvider struct {
	connection *plugin.Connection
}

// Retrieve implements aws.CredentialsProvider. Called by the AWS SDK on every
// signed request, subject to the wrapping CredentialsCache.
func (p *connectionConfigCredentialsProvider) Retrieve(_ context.Context) (creds aws.Credentials, err error) {
	// Belt-and-suspenders: the SDK's Connection.GetConfig acquires an RLock so
	// torn interface reads cannot happen via that path. This recover still
	// converts any other panic inside Retrieve into a clean error the AWS SDK
	// can retry, rather than propagating through the signing middleware.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("connectionConfigCredentialsProvider: panic during Retrieve for connection %q: %v", p.connectionName(), r)
			creds = aws.Credentials{}
		}
	}()

	if p.connection == nil {
		return aws.Credentials{}, errors.New("connectionConfigCredentialsProvider: connection is nil")
	}

	// Read the raw config through the SDK's Connection.GetConfig accessor
	// (which acquires the RLock) and type-assert directly here. Avoid the
	// local aws/connection_config.go GetConfig helper — it normalizes the
	// Regions slice in place and panics on regions = []. Neither belongs in
	// the AWS request signing path: Retrieve only needs the credential fields,
	// and a malformed connection config should not crash signing goroutines
	// deep inside the AWS SDK middleware.
	raw := p.connection.GetConfig()
	cfg, ok := raw.(awsConfig)
	if !ok {
		return aws.Credentials{}, fmt.Errorf("connectionConfigCredentialsProvider: connection %q config is %T, expected awsConfig", p.connection.Name, raw)
	}

	if cfg.AccessKey == nil {
		return aws.Credentials{}, fmt.Errorf("connectionConfigCredentialsProvider: connection %q has no access_key in config", p.connection.Name)
	}
	if cfg.SecretKey == nil {
		return aws.Credentials{}, fmt.Errorf("connectionConfigCredentialsProvider: connection %q has no secret_key in config", p.connection.Name)
	}

	var sessionToken string
	if cfg.SessionToken != nil {
		sessionToken = *cfg.SessionToken
	}

	return aws.Credentials{
		AccessKeyID:     *cfg.AccessKey,
		SecretAccessKey: *cfg.SecretKey,
		SessionToken:    sessionToken,
		Source:          "connectionConfigCredentialsProvider",
		// config.WithCredentialsProvider wraps any provider in a
		// CredentialsCache. The cache will NOT call Retrieve again until the
		// cached value's Expires has passed (or until cache invalidation),
		// so setting Expires too far out defeats the rotation-pickup goal:
		// the cache would hold the original creds in memory long after they
		// were rotated in Connection.Config.
		//
		// 60 seconds matches what the standalone reproduction harness
		// validated to keep rotation latency bounded. Reading
		// Connection.Config is an in-memory type assertion + struct copy,
		// so a short interval is essentially free.
		CanExpire: true,
		Expires:   time.Now().Add(60 * time.Second),
	}, nil
}

func (p *connectionConfigCredentialsProvider) connectionName() string {
	if p.connection == nil {
		return "<nil>"
	}
	return p.connection.Name
}
