package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/backup"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsBackupProtectedResource(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_backup_protected_resource",
		Description: "AWS Backup Protected Resource",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("resource_arn"),
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFoundException", "InvalidParameter", "InvalidParameterValueException"}),
			Hydrate:           getAwsBackupProtectedResource,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsBackupProtectedResources,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "resource_arn",
				Description: "An Amazon Resource Name (ARN) that uniquely identifies a resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_type",
				Description: "The type of Amazon Web Services resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_backup_time",
				Description: "The date and time a resource was last backed up.",
				Type:        proto.ColumnType_TIMESTAMP,
			},

			// Steampipe standard columns
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ResourceArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsBackupProtectedResources(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAwsBackupProtectedResources")

	// Create session
	svc, err := BackupService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &backup.ListProtectedResourcesInput{
		MaxResults: aws.Int64(1000),
	}

	// Limiting the results
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 1 {
				input.MaxResults = types.Int64(1)
			} else {
				input.MaxResults = limit
			}
		}
	}

	err = svc.ListProtectedResourcesPages(
		input,
		func(page *backup.ListProtectedResourcesOutput, lastPage bool) bool {
			for _, resource := range page.Results {
				d.StreamListItem(ctx, resource)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !lastPage
		},
	)
	return nil, err
}

//// HYDRATE FUNCTION

func getAwsBackupProtectedResource(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsBackupProtectedResource")

	// Create session
	svc, err := BackupService(ctx, d)
	if err != nil {
		return nil, err
	}

	var arn string
	if h.Item != nil {
		arn = *h.Item.(*backup.ProtectedResource).ResourceArn
	} else {
		arn = d.KeyColumnQuals["resource_arn"].GetStringValue()
	}

	params := &backup.DescribeProtectedResourceInput{
		ResourceArn: aws.String(arn),
	}

	detail, err := svc.DescribeProtectedResource(params)
	if err != nil {
		plugin.Logger(ctx).Error("getAwsBackupProtectedResource", "DescribeProtectedResource error", err)
		return nil, err
	}

	return detail, nil
}
