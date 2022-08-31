package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/securityhub"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSecurityHub(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_securityhub_hub",
		Description: "AWS Security Hub",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("hub_arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidAccessException", "ResourceNotFoundException"}),
			},
			Hydrate: getSecurityHub,
		},
		List: &plugin.ListConfig{
			Hydrate: listSecurityHubs,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "hub_arn",
				Description: "The ARN of the Hub resource that was retrieved.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "administrator_account",
				Description: "Provides the details for the Security Hub administrator account for the current member account.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSecurityHubAdministratorAccount,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "auto_enable_controls",
				Description: "Whether to automatically enable new controls when they are added to standards that are enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "subscribed_at",
				Description: "The date and time when Security Hub was enabled in the account.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			/// Standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSecurityHubTags,
			},
			{
				Name:        "title",
				Description: "The title of hub. This is a constant value 'default'",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromConstant("default"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HubArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listSecurityHubs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listSecurityHubs")

	// Create session
	svc, err := SecurityHubService(ctx, d)
	if err != nil {
		return nil, err
	}

	// List call
	resp, err := svc.DescribeHub(&securityhub.DescribeHubInput{})

	if err != nil {
		plugin.Logger(ctx).Error("listSecurityHubs", "query_error", err)
		return nil, nil
	}

	d.StreamListItem(ctx, resp)

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSecurityHub(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSecurityHub")

	hubArn := d.KeyColumnQuals["hub_arn"].GetStringValue()

	// get service
	svc, err := SecurityHubService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &securityhub.DescribeHubInput{
		HubArn: &hubArn,
	}

	// Get call
	op, err := svc.DescribeHub(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getSecurityHub", "ERROR", err)
		return nil, err
	}
	return op, nil
}

func getSecurityHubAdministratorAccount(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// get service
	svc, err := SecurityHubService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securityhub_hub.getSecurityHubAdministratorAccount", "service_creation_error", err)
		return nil, err
	}

	// Build the params
	params := &securityhub.GetAdministratorAccountInput{}

	// Get call
	op, err := svc.GetAdministratorAccount(params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securityhub_hub.getSecurityHubAdministratorAccount", "api_error", err)
		return nil, err
	}
	return op.Administrator, nil
}

func getSecurityHubTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSecurityHubTags")

	hubArn := *h.Item.(*securityhub.DescribeHubOutput).HubArn

	// get service
	svc, err := SecurityHubService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &securityhub.ListTagsForResourceInput{
		ResourceArn: &hubArn,
	}

	// Get call
	op, err := svc.ListTagsForResource(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getSecurityHubTags", "ERROR", err)
		return nil, err
	}
	return op, nil
}
