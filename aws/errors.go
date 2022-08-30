package aws

import (
	"context"
	"errors"
	"path"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/smithy-go"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

// isNotFoundError:: function which returns an ErrorPredicate for AWS API calls
func isNotFoundError(notFoundErrors []string) plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		awsConfig := GetConfig(d.Connection)

		// If the get or list hydrate functions have an overriding IgnoreConfig
		// defined using the isNotFoundError function, then it should
		// also check for errors in the "ignore_error_codes" config argument
		allErrors := append(notFoundErrors, awsConfig.IgnoreErrorCodes...)
		if awsErr, ok := err.(awserr.Error); ok {
			// Added to support regex in not found errors
			for _, pattern := range allErrors {
				if ok, _ := path.Match(pattern, awsErr.Code()); ok {
					return true
				}
			}
		}

		var ae smithy.APIError
		if errors.As(err, &ae) {
			// Added to support regex in not found errors
			for _, pattern := range allErrors {
				if ok, _ := path.Match(pattern, ae.ErrorCode()); ok {
					return true
				}
			}
		}
		return false
	}
}

// shouldIgnoreErrorPluginDefault:: Plugin level default function to ignore a set errors for hydrate functions based on "ignore_error_codes" config argument
func shouldIgnoreErrorPluginDefault() plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		if !hasIgnoredErrorCodes(d.Connection) {
			return false
		}

		awsConfig := GetConfig(d.Connection)
		if awsErr, ok := err.(awserr.Error); ok {
			// Added to support regex in ignoring errors
			for _, pattern := range awsConfig.IgnoreErrorCodes {
				if ok, _ := path.Match(pattern, awsErr.Code()); ok {
					return true
				}
			}
		}

		var ae smithy.APIError
		if errors.As(err, &ae) {
			// Added to support regex in not found errors
			for _, pattern := range awsConfig.IgnoreErrorCodes {
				if ok, _ := path.Match(pattern, ae.ErrorCode()); ok {
					return true
				}
			}
		}
		return false
	}
}

func hasIgnoredErrorCodes(connection *plugin.Connection) bool {
	awsConfig := GetConfig(connection)
	return len(awsConfig.IgnoreErrorCodes) > 0
}
