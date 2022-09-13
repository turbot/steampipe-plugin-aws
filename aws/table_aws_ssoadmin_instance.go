package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssoadmin"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsSsoAdminInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssoadmin_instance",
		Description: "AWS SSO Instance",
		List: &plugin.ListConfig{
			Hydrate: listSsoAdminInstances,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The ARN of the SSO instance under which the operation will be executed.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InstanceArn"),
			},
			{
				Name:        "identity_store_id",
				Description: "The identifier of the identity store that is connected to the SSO instance.",
				Type:        proto.ColumnType_STRING,
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InstanceArn"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("InstanceArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listSsoAdminInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listSsoAdminInstances")

	// Create session
	svc, err := SSOAdminService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &ssoadmin.ListInstancesInput{
		MaxResults: aws.Int64(100),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 1 {
				input.MaxResults = aws.Int64(1)
			} else {
				input.MaxResults = limit
			}
		}
	}

	err = svc.ListInstancesPages(
		input,
		func(page *ssoadmin.ListInstancesOutput, isLast bool) bool {
			for _, instance := range page.Instances {
				d.StreamListItem(ctx, instance)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)
	if err != nil {
		plugin.Logger(ctx).Error("listSsoAdminInstances", "ListInstancesPages_error", err)
		return nil, err
	}

	return nil, err
}
