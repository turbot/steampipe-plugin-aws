package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

// isNotFoundErrorWithContext:: function which returns an ErrorPredicate for AWS API calls
func isNotFoundErrorWithContext(IgnoreErrors []string) plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		if awsErr, ok := err.(awserr.Error); ok {

			return helpers.StringSliceContains(IgnoreErrors, awsErr.Code())
		}
		return false
	}
}

func shouldIgnoreErrorTableDefault(IgnoreErrors []string) plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		awsConfig := GetConfig(d.Connection)
		if awsConfig.ShouldIgnoreErrors == nil || (awsConfig.ShouldIgnoreErrors != nil && !*awsConfig.ShouldIgnoreErrors) {
			return false
		}
		if awsErr, ok := err.(awserr.Error); ok {
			// plugin.Logger(ctx).Info("shouldIgnoreErrorTableDefault", "AWS Error CODE", awsErr.Code())
			return helpers.StringSliceContains(IgnoreErrors, awsErr.Code())
		}
		return false
	}
}
