package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/securityhub"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSecurityHub(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_securityhub_hub",
		Description: "AWS Security Hub",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("hub_arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidAccessException", "ResourceNotFoundException"}),
			},
			Hydrate: getSecurityHub,
			Tags:    map[string]string{"service": "securityhub", "action": "DescribeHub"},
		},
		List: &plugin.ListConfig{
			Hydrate: listSecurityHubs,
			Tags:    map[string]string{"service": "securityhub", "action": "DescribeHub"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_SECURITYHUB_SERVICE_ID),
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getSecurityHubAdministratorAccount,
				Tags: map[string]string{"service": "securityhub", "action": "GetAdministratorAccount"},
			},
			{
				Func: getSecurityHubTags,
				Tags: map[string]string{"service": "securityhub", "action": "ListTagsForResource"},
			},
		},
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

			/// Steampie Standard Columns

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

	// Create session
	svc, err := SecurityHubClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securityhub_hub.listSecurityHubs", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// List call
	resp, err := svc.DescribeHub(ctx, &securityhub.DescribeHubInput{})
	if err != nil {
		if strings.Contains(err.Error(), "is not subscribed to AWS Security Hub") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_securityhub_hub.listSecurityHubs", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, resp)

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSecurityHub(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	hubArn := d.EqualsQuals["hub_arn"].GetStringValue()

	// Create session
	svc, err := SecurityHubClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securityhub_hub.getSecurityHub", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &securityhub.DescribeHubInput{
		HubArn: &hubArn,
	}

	// Execute get call
	data, err := svc.DescribeHub(ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), "is not subscribed to AWS Security Hub") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_securityhub_hub.getSecurityHub", "api_error", err)
		return nil, err
	}
	return data, nil
}

func getSecurityHubAdministratorAccount(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := SecurityHubClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securityhub_hub.getSecurityHubAdministratorAccount", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &securityhub.GetAdministratorAccountInput{}

	// Get call
	data, err := svc.GetAdministratorAccount(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securityhub_hub.getSecurityHubAdministratorAccount", "api_error", err)
		return nil, err
	}
	return data.Administrator, nil
}

func getSecurityHubTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	hubArn := *h.Item.(*securityhub.DescribeHubOutput).HubArn

	// Create session
	svc, err := SecurityHubClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securityhub_hub.getSecurityHubTags", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &securityhub.ListTagsForResourceInput{
		ResourceArn: &hubArn,
	}

	// Get call
	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securityhub_hub.getSecurityHubTags", "api_error", err)
		return nil, err
	}
	return op, nil
}
