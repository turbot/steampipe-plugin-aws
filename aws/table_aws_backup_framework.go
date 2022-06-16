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

//// LIST FUNCTION

func listAwsBackupFrameworks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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
