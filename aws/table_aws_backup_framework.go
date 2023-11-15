package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/backup"
	"github.com/aws/aws-sdk-go-v2/service/backup/types"

	backupv1 "github.com/aws/aws-sdk-go/service/backup"

	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// // TABLE DEFINITION
func tableAwsBackupFramework(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_backup_framework",
		Description: "AWS Backup Framework",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("framework_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValueException"}),
			},
			Hydrate: getAwsBackupFramework,
			Tags:    map[string]string{"service": "backup", "action": "DescribeFramework"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsBackupFrameworks,
			Tags:    map[string]string{"service": "backup", "action": "ListFrameworks"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsBackupFramework,
				Tags: map[string]string{"service": "backup", "action": "DescribeFramework"},
			},
			{
				Func: listAwsBackupFrameworkTags,
				Tags: map[string]string{"service": "backup", "action": "ListTags"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(backupv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "framework_name",
				Description: "The unique name of a backup framework.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "An Amazon Resource Name (ARN) that uniquely identifies a backup framework resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FrameworkArn"),
			},
			{
				Name:        "framework_description",
				Description: "An optional description of the backup framework.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "deployment_status",
				Description: "The deployment status of a backup framework.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The date and time that a framework was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "number_of_controls",
				Description: "The number of controls contained by the framework.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getNumberOfControls,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "framework_status",
				Description: "The framework status based on recording statuses for resources governed by the framework (ACTIVE | PARTIALLY_ACTIVE | INACTIVE | UNAVAILABLE).",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsBackupFramework,
			},
			{
				Name:        "framework_controls",
				Description: "A list of the controls that make up the framework. Each control in the list has a name, input parameters, and scope.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsBackupFramework,
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAwsBackupFrameworkTags,
				Transform:   transform.FromValue(),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FrameworkName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("FrameworkArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

func getNumberOfControls(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var value int

	switch item := h.Item.(type) {
	case *backup.DescribeFrameworkOutput:
		value = len(item.FrameworkControls)
	case types.Framework:
		value = int(item.NumberOfControls)
	}

	return value, nil
}

//// LIST FUNCTION

func listAwsBackupFrameworks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := BackupClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_framework.listAwsBackupFrameworks", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
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

	input := &backup.ListFrameworksInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := backup.NewListFrameworksPaginator(svc, input, func(o *backup.ListFrameworksPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_backup_framework.listAwsBackupFrameworks", "api_error", err)
			return nil, err
		}

		for _, items := range output.Frameworks {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsBackupFramework(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// AWS Backup service is available in all regions. However, the AWS Backup audit manager, which is newly introduced under the Backup service, is not supported in all regions.
	// Due to this reason, we could not put a check based on the service endpoint and had to check the region code directly.
	// https://aws.amazon.com/about-aws/whats-new/2022/05/aws-backup-audit-manager-adds-amazon-s3-storage-gateway/#:~:text=AWS%20Backup%20Audit%20Manager%20is,Middle%20East%20(Bahrain)%20Regions.
	region := d.EqualsQualString(matrixKeyRegion)
	if region == "ap-northeast-3" {
		return nil, nil
	}

	// Create Session
	svc, err := BackupClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_framework.getAwsBackupFramework", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	var name string
	if h.Item != nil {
		framework := h.Item.(types.Framework)
		name = *framework.FrameworkName
	} else {
		name = d.EqualsQuals["framework_name"].GetStringValue()
	}

	// check if id is empty
	if name == "" {
		return nil, nil
	}

	params := &backup.DescribeFrameworkInput{
		FrameworkName: aws.String(name),
	}

	op, err := svc.DescribeFramework(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "ResourceNotFoundException" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_backup_framework.getAwsBackupFramework", "api_error", err)
		return nil, err
	}

	return op, nil
}

func listAwsBackupFrameworkTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := BackupClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_framework.listAwsBackupFrameworkTags", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	var arn *string

	switch item := h.Item.(type) {
	case types.Framework:
		arn = item.FrameworkArn
	case *backup.DescribeFrameworkOutput:
		arn = item.FrameworkArn
	}

	// Build the params
	params := backup.ListTagsInput{
		ResourceArn: aws.String(*arn),
		MaxResults:  aws.Int32(1000),
	}

	paginator := backup.NewListTagsPaginator(svc, &params, func(o *backup.ListTagsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	tags := make(map[string]string)

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_api_gateway_rest_api.listRestAPI", "api_error", err)
			return nil, err
		}

		for k, v := range output.Tags {
			tags[k] = v
		}

	}
	return tags, nil
}
