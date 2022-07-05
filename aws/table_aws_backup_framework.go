package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/backup"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION
func tableAwsBackupFramework(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_backup_framework",
		Description: "AWS Backup Framework",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("framework_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidParameterValueException"}),
			},
			Hydrate: getAwsBackupFramework,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsBackupFrameworks,
		},
		GetMatrixItem: BuildRegionList,
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
	case *backup.Framework:
		value = int(*item.NumberOfControls)
	}

	return value, nil
}

//// LIST FUNCTION

func listAwsBackupFrameworks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// AWS Backup service is available in all regions. However, the AWS Backup audit manager, which is newly introduced under the Backup service, is not supported in all regions.
	// Due to this reason, we could not put a check based on the service endpoint and had to check the region code directly.
	// https://aws.amazon.com/about-aws/whats-new/2022/05/aws-backup-audit-manager-adds-amazon-s3-storage-gateway/#:~:text=AWS%20Backup%20Audit%20Manager%20is,Middle%20East%20(Bahrain)%20Regions.
	region := d.KeyColumnQualString(matrixKeyRegion)
	if region == "ap-northeast-3" {
		return nil, nil
	}

	// Create session
	svc, err := BackupService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &backup.ListFrameworksInput{
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

	err = svc.ListFrameworksPages(
		input,
		func(output *backup.ListFrameworksOutput, lastPage bool) bool {
			for _, plan := range output.Frameworks {
				d.StreamListItem(ctx, plan)

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

//// HYDRATE FUNCTIONS

func getAwsBackupFramework(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// AWS Backup service is available in all regions. However, the AWS Backup audit manager, which is newly introduced under the Backup service, is not supported in all regions.
	// Due to this reason, we could not put a check based on the service endpoint and had to check the region code directly.
	// https://aws.amazon.com/about-aws/whats-new/2022/05/aws-backup-audit-manager-adds-amazon-s3-storage-gateway/#:~:text=AWS%20Backup%20Audit%20Manager%20is,Middle%20East%20(Bahrain)%20Regions.
	region := d.KeyColumnQualString(matrixKeyRegion)
	if region == "ap-northeast-3" {
		return nil, nil
	}

	// Create Session
	svc, err := BackupService(ctx, d)
	if err != nil {
		return nil, err
	}

	var name string
	if h.Item != nil {
		framework := h.Item.(*backup.Framework)
		name = *framework.FrameworkName
	} else {
		name = d.KeyColumnQuals["framework_name"].GetStringValue()
	}

	// check if id is empty
	if name == "" {
		return nil, nil
	}

	params := &backup.DescribeFrameworkInput{
		FrameworkName: aws.String(name),
	}

	op, err := svc.DescribeFramework(params)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "ResourceNotFoundException" {
				return nil, nil
			}
		}
		return nil, err
	}

	return op, nil
}

func listAwsBackupFrameworkTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := BackupService(ctx, d)
	if err != nil {
		return nil, err
	}

	var arn *string

	switch item := h.Item.(type) {
	case *backup.Framework:
		arn = item.FrameworkArn
	case *backup.DescribeFrameworkOutput:
		arn = item.FrameworkArn
	}

	// Build the params
	params := backup.ListTagsInput{
		ResourceArn: aws.String(*arn),
		MaxResults:  aws.Int64(1000),
	}

	tags := make(map[string]string)
	pagesLeft := true
	for pagesLeft {
		keyTags, err := svc.ListTags(&params)
		if err != nil {
			plugin.Logger(ctx).Error("listAwsBackupFrameworkTags", "ListTags_error", err)
			return nil, err
		}

		for k, v := range keyTags.Tags {
			tags[k] = *v
		}

		if keyTags.NextToken != nil {
			params.NextToken = keyTags.NextToken
		} else {
			pagesLeft = false
		}
	}

	return tags, nil
}
