package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/backup"
	"github.com/aws/smithy-go"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func getAwsBackupResourceTags(ctx context.Context, d *plugin.QueryData, arn string) (interface{}, error) {
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
	plugin.Logger(ctx).Debug("backup_tags.getAwsBackupResourceTags", "ListTagsOutput", op)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			plugin.Logger(ctx).Debug("backup_tags.getAwsBackupResourceTags", "smithy.APIError", ae)
			if ae.ErrorCode() == "ResourceNotFoundException" {
				return &backup.ListTagsOutput{
					Tags: map[string]string{},
				}, nil
			}
		}
		plugin.Logger(ctx).Error("backup_tags.getAwsBackupResourceTags", "api_error", err)
		return nil, err
	}

	if op.Tags == nil {
		op.Tags = map[string]string{}
	}

	return op, nil
}
