package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/securityhub"
	"github.com/aws/aws-sdk-go-v2/service/securityhub/types"
	"github.com/aws/smithy-go"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSecurityHubStandardsSubscription(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_securityhub_standards_subscription",
		Description: "AWS Security Hub Standards Subscription",
		List: &plugin.ListConfig{
			Hydrate: listSecurityHubStandardsSubcriptions,
		},
		GetMatrixItemFunc: BuildRegionList,
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
				Name:        "standards_status_reason_code",
				Description: "The reason code that represents the reason for the current status of a standard subscription.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     GetEnabledStandards,
				Transform:   transform.FromField("StandardsStatusReason.StatusReasonCode"),
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

func listSecurityHubStandardsSubcriptions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := SecurityHubClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securityhub_standards_subscription.listSecurityHubStandardsSubcriptions", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(100)
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

	input := &securityhub.DescribeStandardsInput{
		MaxResults: maxLimit,
	}

	paginator := securityhub.NewDescribeStandardsPaginator(svc, input, func(o *securityhub.DescribeStandardsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_securityhub_standards_subscription.listSecurityHubStandardsSubcriptions", "api_error", err)
			return nil, err
		}

		for _, standards := range output.Standards {
			d.StreamListItem(ctx, standards)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func GetEnabledStandards(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	standardArn := *h.Item.(types.Standard).StandardsArn

	// Create session
	svc, err := SecurityHubClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securityhub_standards_subscription.GetEnabledStandards", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the input
	input := &securityhub.GetEnabledStandardsInput{}

	// Get call
	standardsSubscriptions, err := svc.GetEnabledStandards(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securityhub_standards_subscription.GetEnabledStandards", "api_error", err)
		var ae smithy.APIError
		if errors.As(err, &ae) {
			// If the service is not enabled, API throws InvalidAccessException error
			if ae.ErrorCode() == "InvalidAccessException" {
				return nil, nil
			}
		}
	}

	for _, item := range standardsSubscriptions.StandardsSubscriptions {
		if *item.StandardsArn == standardArn {
			return item, nil
		}
	}
	return nil, nil
}
