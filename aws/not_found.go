package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

// function which returns an ErrorPredicate for AWS API calls
func isNotFoundError(notFoundErrors []string) plugin.ErrorPredicate {
	return func(err error) bool {
		if awsErr, ok := err.(awserr.Error); ok {
			return helpers.StringSliceContains(notFoundErrors, awsErr.Code())
		}
		return false
	}
}

func ignoreAccessDeniedError(ctx context.Context, AccessDeniedErrors []string) plugin.ErrorPredicate {
	// if
	return func(err error) bool {
		if awsErr, ok := err.(awserr.Error); ok {
			// plugin.Logger(ctx).Info("shouldIgnoreError", "AWS Error CODE", awsErr.Code())
			return helpers.StringSliceContains(AccessDeniedErrors, awsErr.Code())
		}
		return false
	}
}

func shouldIgnoreErrorTableDefault(IgnoreErrors []string) plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		awsConfig := GetConfig(d.Connection)
		if awsConfig.IgnoreAccessDeniedErrors == nil || (awsConfig.IgnoreAccessDeniedErrors != nil && !*awsConfig.IgnoreAccessDeniedErrors) {
			return false
		}
		if awsErr, ok := err.(awserr.Error); ok {
			plugin.Logger(ctx).Info("shouldIgnoreErrorTableDefault", "AWS Error CODE", awsErr.Code())
			return helpers.StringSliceContains(IgnoreErrors, awsErr.Code())
		}
		return false
	}
}
func isNotFoundErrorWithContext(IgnoreErrors []string) plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		if awsErr, ok := err.(awserr.Error); ok {
			// plugin.Logger(ctx).Info("isNotFoundErrorWithContext", "AWS Error CODE", awsErr.Code())
			return helpers.StringSliceContains(IgnoreErrors, awsErr.Code())
		}
		return false
	}
}
