package aws

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
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
