package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/securityhub"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsSecurityHubStandardSubscription(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_securityhub_standards_subscription",
		Description: "AWS Securityhub standards",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidAccessException"}),
			Hydrate:           getSecurityHubStandardSubcription,
		},
		List: &plugin.ListConfig{
			Hydrate: listSecurityHubStandardSubcriptions,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the standard.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "standards_arn",
				Description: "The ARN of a standard.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The description of the standard.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enabled_by_default",
				Description: "Indicates whether the standard is enabled by default.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "standards_status",
				Description: "The status of the standard subscription.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     GetEnabledStandards,
			},
			{
				Name:        "standards_subscription_arn",
				Description: "The ARN of a resource that represents your subscription to a supported standard.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     GetEnabledStandards,
			},
			// JSON columns
			{
				Name:        "standards_input",
				Description: "A key-value pair of input for the standard.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     GetEnabledStandards,
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("StandardsArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listSecurityHubStandardSubcriptions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listSecurityHubStandardSubcriptions", "AWS_REGION", region)

	// Create session
	svc, err := SecurityHubService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	resp, err := svc.DescribeStandards(&securityhub.DescribeStandardsInput{})
	if err != nil {
		plugin.Logger(ctx).Error("listSecurityHubStandardSubcriptions", "query_error", err)
		return nil, nil
	}
	for _, standards := range resp.Standards {
		d.StreamListItem(ctx, standards)
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSecurityHubStandardSubcription(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSecurityHubStandardSubcription")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	name := d.KeyColumnQuals["name"].GetStringValue()

	// get service
	svc, err := SecurityHubService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the input
	input := &securityhub.DescribeStandardsInput{}

	// Get call
	standards, err := svc.DescribeStandards(input)
	if err != nil {
		return nil, err
	}

	for _, item := range standards.Standards {
		if *item.Name == name {
			return item, nil
		}
	}
	return nil, nil
}

func GetEnabledStandards(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("GetEnabledStandards")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	standardArn := *h.Item.(*securityhub.Standard).StandardsArn
	// get service
	svc, err := SecurityHubService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the input
	input := &securityhub.GetEnabledStandardsInput{}

	// Get call
	standardsSubscriptions, err := svc.GetEnabledStandards(input)
	if err != nil {
		return nil, err
	}

	for _, item := range standardsSubscriptions.StandardsSubscriptions {
		if *item.StandardsArn == standardArn {
			return item, nil
		}
	}
	return nil, nil
}
