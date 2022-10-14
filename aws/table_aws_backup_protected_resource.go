package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/backup"
	"github.com/aws/aws-sdk-go-v2/service/backup/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsBackupProtectedResource(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_backup_protected_resource",
		Description: "AWS Backup Protected Resource",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("resource_arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundErrorV2([]string{"ResourceNotFoundException", "InvalidParameter", "InvalidParameterValueException"}),
			},
			Hydrate: getAwsBackupProtectedResource,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsBackupProtectedResources,
		},
		GetMatrixItemFunc: BuildRegionList,
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

	// Create session
	svc, err := BackupClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_protected_resource.listAwsBackupProtectedResources", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// unsupported region check
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	input := &backup.ListProtectedResourcesInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := backup.NewListProtectedResourcesPaginator(svc, input, func(o *backup.ListProtectedResourcesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_backup_protected_resource.listAwsBackupProtectedResources", "api_error", err)
			return nil, err
		}

		for _, items := range output.Results {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTION

func getAwsBackupProtectedResource(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := BackupClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_protected_resource.getAwsBackupProtectedResource", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// unsupported region check
		return nil, nil
	}

	var arn string
	if h.Item != nil {
		arn = *h.Item.(types.ProtectedResource).ResourceArn
	} else {
		arn = d.KeyColumnQuals["resource_arn"].GetStringValue()
	}

	params := &backup.DescribeProtectedResourceInput{
		ResourceArn: aws.String(arn),
	}

	detail, err := svc.DescribeProtectedResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_protected_resource.getAwsBackupProtectedResource", "api_error", err)
		return nil, err
	}

	return detail, nil
}
