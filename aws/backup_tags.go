package aws

import (
	"context"
	"errors"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/backup"
	"github.com/aws/smithy-go"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func getAwsBackupResourceTags(ctx context.Context, d *plugin.QueryData, arn string, pattern string) (interface{}, error) { // Create a regular expression object
	re := regexp.MustCompile(pattern)

	// Only return the tags associated with the resovery point
	if !re.MatchString(arn) {
		return nil, nil
	}

	// Create Session
	svc, err := BackupClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("backup_tags.getAwsBackupResourceTags", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &backup.ListTagsInput{
		ResourceArn: aws.String(arn),
	}

	op, err := svc.ListTags(ctx, params)
	if err != nil {

		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "ResourceNotFoundException" {
				return &backup.ListTagsOutput{}, nil
			}
		}
		plugin.Logger(ctx).Error("backup_tags.getAwsBackupResourceTags", "api_error", err)
		return nil, err
	}

	if op.Tags == nil {
		return nil, nil
	}

	return op, nil
}
