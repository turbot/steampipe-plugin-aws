package aws

import (
	"context"
	"path"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	// "github.com/turbot/go-kit/helpers"
)

// isNotFoundError:: function which returns an ErrorPredicate for AWS API calls
func isNotFoundError(notFoundErrors []string) plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		awsConfig := GetConfig(d.Connection)
		// If the get or list hydrate functions have a overriding IgnoreConfig defined using isNotFoundError function - It should also check for errors in "ignored_error_codes"
		allErrors := append(notFoundErrors, awsConfig.IgnoredErrorCodes...)
		if awsErr, ok := err.(awserr.Error); ok {
			plugin.Logger(ctx).Info("isNotFoundError", "AWS Error CODE", awsErr.Code())
			// return helpers.StringSliceContains(allErrors, awsErr.Code())

			// Added to support regex in not found errors
			for _, pattern := range allErrors {
				if ok, _ := path.Match(pattern, awsErr.Code()); ok {
					return true
				}
			}
		}
		return false
	}
}

// shouldIgnoreErrorPluginDefault:: Plugin level default function to ignore a set errors for hydrate functions based on `ignored_error_codes`
func shouldIgnoreErrorPluginDefault() plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		if !hasIgnoredErrorCodes(d.Connection) {
			return false
		}
		awsConfig := GetConfig(d.Connection)
		if awsErr, ok := err.(awserr.Error); ok {
			plugin.Logger(ctx).Info("shouldIgnoreErrorPluginDefault", "AWS Error CODE", awsErr.Code())
			// return helpers.StringSliceContains(awsConfig.IgnoredErrorCodes, awsErr.Code())

			// Added to support regex in ignoring errors
			for _, pattern := range awsConfig.IgnoredErrorCodes {
				if ok, _ := path.Match(pattern, awsErr.Code()); ok {
					return true
				}
			}
		}
		return false
	}
}

func hasIgnoredErrorCodes(connection *plugin.Connection) bool {
	awsConfig := GetConfig(connection)
	return len(awsConfig.IgnoredErrorCodes) > 0
}
