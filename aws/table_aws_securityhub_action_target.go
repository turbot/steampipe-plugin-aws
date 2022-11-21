package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/securityhub"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
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
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidInputException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listSecurityHubActionTargets,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
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

	// Create session
	svc, err := SecurityHubClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securityhub_action_target.go.listSecurityHubActionTargets", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
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

	// Filter parameter is not supported yet in this SDK version so optional quals can not be implemented
	input := &securityhub.DescribeActionTargetsInput{
		MaxResults: maxLimit,
	}

	// List call
	paginator := securityhub.NewDescribeActionTargetsPaginator(svc, input, func(o *securityhub.DescribeActionTargetsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			if strings.Contains(err.Error(), "not subscribed") {
				return nil, nil
			}
			plugin.Logger(ctx).Error("aws_securityhub_action_target.go.listSecurityHubActionTargets", "api_error", err)
			return nil, err
		}

		for _, actionTarget := range output.ActionTargets {
			d.StreamListItem(ctx, actionTarget)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSecurityHubActionTarget(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	arn := d.KeyColumnQuals["arn"].GetStringValue()

	// Create session
	svc, err := SecurityHubClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securityhub_action_target.go.getSecurityHubActionTarget", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &securityhub.DescribeActionTargetsInput{
		ActionTargetArns: []string{arn},
	}

	// Get call
	data, err := svc.DescribeActionTargets(ctx, params)
	if err != nil {
		// Handle unsupported and inactive region exceptions
		if strings.Contains(err.Error(), "not subscribed") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_securityhub_action_target.getSecurityHubActionTarget", "api_error", err)
		return nil, err
	}

	if len(data.ActionTargets) > 0 {
		return data.ActionTargets[0], nil
	}

	return nil, nil
}
