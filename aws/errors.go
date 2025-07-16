package aws

import (
	"context"
	"errors"
	"path"
	"regexp"

	"github.com/aws/smithy-go"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// shouldIgnoreErrors:: function which returns an ErrorPredicate for AWS API calls
func shouldIgnoreErrors(notFoundErrors []string) plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		awsConfig := GetConfig(d.Connection)

		// If the get or list hydrate functions have an overriding IgnoreConfig
		// defined using the shouldIgnoreErrors function, then it should
		// also check for errors in the "ignore_error_codes" config argument
		allErrors := append(notFoundErrors, awsConfig.IgnoreErrorCodes...)
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

// shouldIgnoreErrorPluginDefault:: Plugin level default function to ignore a set errors for hydrate functions based on "ignore_error_codes" and "ignore_error_messages" config argument
func shouldIgnoreErrorPluginDefault() plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		if !hasIgnoredErrorCodesOrMessages(d.Connection) {
			return false
		}

		awsConfig := GetConfig(d.Connection)

		logger := plugin.Logger(ctx)

		// Add to support regex match as per error message
		for _, pattern := range awsConfig.IgnoreErrorMessages {
			// Validate regex pattern
			re, er := regexp.Compile(pattern)
			if er != nil {
				panic(er.Error() + " the regex pattern configured in 'ignore_error_messages' is invalid. Edit your connection configuration file and then restart Steampipe")
			}
			result := re.MatchString(err.Error())
			if result {
				logger.Debug("errors.shouldIgnoreErrors", "ignore_error_message", err.Error())
				return true
			}
		}

		var ae smithy.APIError
		if errors.As(err, &ae) {
			// Added to support regex in not found errors
			for _, pattern := range awsConfig.IgnoreErrorCodes {
				if ok, _ := path.Match(pattern, ae.ErrorCode()); ok {
					logger.Debug("errors.shouldIgnoreErrorPluginDefault", "ignore_error_code", err.Error())
					return true
				}
			}
		}
		return false
	}
}

func hasIgnoredErrorCodesOrMessages(connection *plugin.Connection) bool {
	awsConfig := GetConfig(connection)
	return len(awsConfig.IgnoreErrorCodes) > 0 || len(awsConfig.IgnoreErrorMessages) > 0
}
