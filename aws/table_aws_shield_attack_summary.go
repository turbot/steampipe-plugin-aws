package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/shield"
	"github.com/aws/aws-sdk-go-v2/service/shield/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsShieldAttackSummary(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_shield_attack_summary",
		Description: "AWS Shield Attack Summary",
		List: &plugin.ListConfig{
			Hydrate: listAwsShieldAttackSummaries,
			KeyColumns: plugin.OptionalColumns([]string{"resource_arn", "start_time", "end_time"}),
			Tags:    map[string]string{"service": "shield", "action": "ListAttacks"},
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "attack_id",
				Description: "The unique identifier (ID) of the attack.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AttackId"),
			},
			{
				Name:        "resource_arn",
				Description: "The ARN (Amazon Resource Name) of the Amazon Web Services resource that was attacked.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceArn"),
			},
			{
				Name:        "start_time",
				Description: "The start time of the attack.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("StartTime"),
			},
			{
				Name:        "end_time",
				Description: "The end time of the attack.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("EndTime"),
			},
			{
				Name:        "attack_vectors",
				Description: "The list of attacks for the time period.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AttackVectors"),
			},
			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AttackId"),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listAwsShieldAttackSummaries(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := ShieldClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_shield_attack_summary.listAwsShieldAttackSummaries", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	queryResultLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		queryResultLimit = min(queryResultLimit, int32(*d.QueryContext.Limit))
	}

	input := &shield.ListAttacksInput{
		MaxResults: aws.Int32(queryResultLimit),
	}

	if d.Quals["resource_arn"] != nil {
		for _, q := range d.Quals["resource_arn"].Quals {
			input.ResourceArns = []string{}
			input.ResourceArns = append(input.ResourceArns, q.Value.GetStringValue())
		}
	}

	if d.Quals["start_time"] != nil {
		for _, q := range d.Quals["start_time"].Quals {
			timestamp := q.Value.GetTimestampValue().AsTime()
			input.StartTime = &types.TimeRange{}

			input.StartTime.FromInclusive = aws.Time(timestamp)
			input.StartTime.ToExclusive = aws.Time(timestamp)
		}
	}

	if d.Quals["end_time"] != nil {
		for _, q := range d.Quals["end_time"].Quals {
			timestamp := q.Value.GetTimestampValue().AsTime()
			input.EndTime = &types.TimeRange{}

			input.EndTime.FromInclusive = aws.Time(timestamp)
			input.EndTime.ToExclusive = aws.Time(timestamp)
		}
	}

	paginator := shield.NewListAttacksPaginator(svc, input, func(o *shield.ListAttacksPaginatorOptions) {
		o.Limit = queryResultLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_shield_attack_summary.listAwsShieldAttackSummaries", "api_error", err)
			return nil, err
		}

		for _, items := range output.AttackSummaries {
			d.StreamListItem(ctx, &items)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}