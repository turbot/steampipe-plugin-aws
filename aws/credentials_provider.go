package aws

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// connectionConfigCredentialsProvider is an aws.CredentialsProvider that reads
// access_key / secret_key / session_token from the steampipe connection config
// on every Retrieve() call, instead of caching them at construction time.
//
// This matters when the connection's credentials are rotated mid-query (e.g.
// Turbot Pipes temp-credential refresh). The steampipe plugin SDK mutates
// Connection.Config in-place when new config arrives via UpdateConnectionConfigs
// (see steampipe-plugin-sdk/plugin/plugin_connection_config.go upsertConnectionData,
// where d.Connection.Config = configStruct is the live mutation). A
// credentials.NewStaticCredentialsProvider built once at aws.Config construction
// would freeze the original token values, and a goroutine holding that
// aws.Config would keep signing requests with the original (eventually expired)
// token forever — surfacing as ExpiredToken errors in long-running queries
// or in queries that started before a credential rotation reached the pod.
//
// By re-reading Connection.Config on every Retrieve(), in-flight goroutines
// holding the same aws.Config automatically pick up rotated credentials on the
// next signing operation.
type connectionConfigCredentialsProvider struct {
	connection *plugin.Connection
}

// Retrieve implements aws.CredentialsProvider. It is called by the AWS SDK on
// every signed request, subject to the SDK's CredentialsCache (which we hint at
// via the short Expires below to force re-reads well before the underlying
// STS token TTL).
func (p *connectionConfigCredentialsProvider) Retrieve(_ context.Context) (aws.Credentials, error) {
	if p.connection == nil {
		return aws.Credentials{}, errors.New("connectionConfigCredentialsProvider: connection is nil")
	}
	cfg := GetConfig(p.connection)
	if cfg.AccessKey == nil || cfg.SecretKey == nil {
		return aws.Credentials{}, errors.New("connectionConfigCredentialsProvider: access_key or secret_key missing from connection config")
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
		// Hint the SDK's CredentialsCache to re-call Retrieve well before the
		// underlying STS token TTL. The maximum STS TTL for chained AssumeRole
		// is 1 hour, so 50 minutes leaves 10 minutes of safety margin for
		// rotation propagation.
		CanExpire: true,
		Expires:   time.Now().Add(50 * time.Minute),
	}, nil
}
