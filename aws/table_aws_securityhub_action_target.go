package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/securityhub"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
)

//// TABLE DEFINITION

func tableAwsSecurityHubActionTarget(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_securityhub_action_target",
		Description: "AWS Security Hub Action Target",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			Hydrate:    getSecurityHubActionTarget,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidInputException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listSecurityHubActionTargets,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the action target.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The ARN for the target action.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ActionTargetArn"),
			},
			{
				Name:        "description",
				Description: "The description of the target action.",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe standard columns
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
				Transform:   transform.FromField("ActionTargetArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listSecurityHubActionTargets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listSecurityHubActionTargets")

	// Create session
	svc, err := SecurityHubService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Filter parameter is not supported yet in this SDK version so optional quals can not be implemented
	input := &securityhub.DescribeActionTargetsInput{
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

	// List call
	err = svc.DescribeActionTargetsPages(
		input,
		func(page *securityhub.DescribeActionTargetsOutput, isLast bool) bool {
			for _, actionTarget := range page.ActionTargets {
				d.StreamListItem(ctx, actionTarget)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	if err != nil {
		// Handle unsupported and inactive region exceptions
		if strings.Contains(err.Error(), "not subscribed") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("listSecurityHubActionTargets", "list", err)
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getSecurityHubActionTarget(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getSecurityHubActionTarget")

	arn := d.KeyColumnQuals["arn"].GetStringValue()

	// Create session
	svc, err := SecurityHubService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &securityhub.DescribeActionTargetsInput{
		ActionTargetArns: aws.StringSlice([]string{arn}),
	}

	// Get call
	data, err := svc.DescribeActionTargets(params)
	if err != nil {
		// Handle unsupported and inactive region exceptions
		if strings.Contains(err.Error(), "not subscribed") {
			return nil, nil
		}
		logger.Error("getSecurityHubActionTarget", "get", err)
		return nil, err
	}

	if len(data.ActionTargets) > 0 {
		return data.ActionTargets[0], nil
	}

	return nil, nil
}
