package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

// isNotFoundErrorWithContext:: function which returns an ErrorPredicate for AWS API calls
func isNotFoundErrorWithContext(notFoundErrors []string) plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		awsConfig := GetConfig(d.Connection)
		allErrors := notFoundErrors
		// If the get or list hydrate functions have a separate IgnoreConfig defined using isNotFoundErrorWithContext function - It will also check for "should_ignore_flag"
		if awsConfig.ShouldIgnoreErrors != nil && *awsConfig.ShouldIgnoreErrors {
			allErrors = append(allErrors, accessDeniedErrors...)
		}
		if awsErr, ok := err.(awserr.Error); ok {
			return helpers.StringSliceContains(allErrors, awsErr.Code())
		}
		return false
	}
}
