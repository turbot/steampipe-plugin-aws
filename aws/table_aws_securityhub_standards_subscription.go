package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/securityhub"
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
	plugin.Logger(ctx).Trace("listSecurityHubStandardsSubcriptions")

	// Create session
	svc, err := SecurityHubService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &securityhub.DescribeStandardsInput{
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

	err = svc.DescribeStandardsPages(
		input,
		func(page *securityhub.DescribeStandardsOutput, isLast bool) bool {
			for _, standards := range page.Standards {
				d.StreamListItem(ctx, standards)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func GetEnabledStandards(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("GetEnabledStandards")

	standardArn := *h.Item.(*securityhub.Standard).StandardsArn
	// get service
	svc, err := SecurityHubService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the input
	input := &securityhub.GetEnabledStandardsInput{}

	// Get call
	standardsSubscriptions, err := svc.GetEnabledStandards(input)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "InvalidAccessException" {
				return nil, nil
			}
			return nil, err
		}
	}

	for _, item := range standardsSubscriptions.StandardsSubscriptions {
		if *item.StandardsArn == standardArn {
			return item, nil
		}
	}
	return nil, err
}
