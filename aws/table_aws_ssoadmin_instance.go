package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/ssoadmin"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsSsoAdminInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssoadmin_instance",
		Description: "AWS SSO Instance",
		List: &plugin.ListConfig{
			Hydrate: listSsoAdminInstances,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "identity_store_id",
				Description: "The identifier of the identity store that is connected to the SSO instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_arn",
				Description: "The ARN of the SSO instance under which the operation will be executed.",
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
	// Create session
	svc, err := SSOAdminService(ctx, d)
	if err != nil {
		return nil, err
	}

	err = svc.ListInstancesPages(
		&ssoadmin.ListInstancesInput{},
		func(page *ssoadmin.ListInstancesOutput, isLast bool) bool {
			for _, instance := range page.Instances {
				d.StreamListItem(ctx, instance)
			}
			return !isLast
		},
	)

	return nil, err
}
