package aws

import (
	"context"
	"slices"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/shield"
	"github.com/aws/aws-sdk-go-v2/service/shield/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/quals"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsShieldAttack(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_shield_attack",
		Description: "AWS Shield Attack",
		Get: &plugin.GetConfig{
			Hydrate:	getAttack,
			KeyColumns: plugin.SingleColumn("attack_id"),
			Tags:		map[string]string{"service": "shield", "action": "DescribeAttack"},
		},
		List: &plugin.ListConfig{
			Hydrate: 	listAttacks,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "resource_arn",
					Require: plugin.Optional,
					Operators: []string{"="},
				},
				{
					Name:    "start_time",
					Require: plugin.Optional,
					Operators: []string{"=", ">", ">=", "<", "<="},
				},
				{
					Name:    "end_time",
					Require: plugin.Optional,
					Operators: []string{"=", ">", ">=", "<", "<="},
				},
			},
			Tags:    	map[string]string{"service": "shield", "action": "ListAttacks"},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "attack_id",
				Description: "The unique identifier (ID) of the attack.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_arn",
				Description: "The ARN (Amazon Resource Name) of the Amazon Web Services resource that was attacked.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "start_time",
				Description: "The start time of the attack.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "end_time",
				Description: "The end time of the attack.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "attack_vectors",
				Description: "The list of attacks for the time period.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "sub_resources",
				Description: "If applicable, additional detail about the resource being attacked, for example, IP address or URL.",
				Hydrate:	 getAttack,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "attack_counters",
				Description: "List of counters that describe the attack for the specified time period.",
				Hydrate:	 getAttack,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "attack_properties",
				Description: "The array of objects that provide details of the Shield event.",
				Hydrate:	 getAttack,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "mitigations",
				Description: "List of mitigation actions taken for the attack.",
				Hydrate:	 getAttack,
				Type:        proto.ColumnType_JSON,
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

//// HELPER FUNCTIONS

type AttackExtended struct {
	AttackVectors []types.AttackVectorDescription
	AttackCounters []types.SummarizedCounter
	AttackId *string
	AttackProperties []types.AttackProperty
	EndTime *time.Time
	Mitigations []types.Mitigation
	ResourceArn *string
	StartTime *time.Time
	SubResources []types.SubResourceSummary
}

func getAttackVectors(attack shield.DescribeAttackOutput) ([]types.AttackVectorDescription) {
	attackVectors := []types.AttackVectorDescription{}

	if attack.Attack.SubResources == nil {
		return attackVectors
	}

	for _, subResource := range attack.Attack.SubResources {
		if subResource.AttackVectors == nil {
			continue
		}
		for _, attackVector := range subResource.AttackVectors {
			attackVectors = append(attackVectors, types.AttackVectorDescription{
				VectorType: attackVector.VectorType,
			})
		}
	}

	return attackVectors
}

func getTimeRange(quals quals.QualSlice) *types.TimeRange {
	fromInclusives := []time.Time{}
	toExclusives := []time.Time{}

	for _, q := range quals {
		operator := q.Operator
		timestamp := q.Value.GetTimestampValue().AsTime()

		var fromInclusive time.Time
		var toExclusive time.Time

		// https://aws.amazon.com/about-aws/whats-new/2016/12/introducing-aws-shield/
		introductionDateOfShield := time.Date(2016, time.December, 1, 0, 0, 0, 0, time.UTC)

		switch operator {
			case "=":
				return &types.TimeRange{
					FromInclusive: aws.Time(timestamp),
					ToExclusive: aws.Time(timestamp.Add(time.Second * 1)),
				}
			case ">":
				fromInclusive = timestamp.Add(time.Second * 1)
				toExclusive = time.Now()
			case ">=":
				fromInclusive = timestamp
				toExclusive = time.Now()
			case "<":
				fromInclusive = introductionDateOfShield
				toExclusive = timestamp
			case "<=":
				fromInclusive = introductionDateOfShield
				toExclusive = timestamp.Add(time.Second * 1)
		}

		fromInclusives = append(fromInclusives, fromInclusive)
		toExclusives = append(toExclusives, toExclusive)
	}

	return &types.TimeRange{
		FromInclusive: aws.Time(slices.MaxFunc(fromInclusives, time.Time.Compare)),
		ToExclusive: aws.Time(slices.MinFunc(toExclusives, time.Time.Compare)),
	}
}

//// HYDRATE FUNCTIONS

func getAttack(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var attackId string
	if h.Item != nil {
		attackId = *h.Item.(*types.AttackSummary).AttackId
	} else {
		attackId = d.EqualsQualString("attack_id")
	}

	if attackId == "" {
		return nil, nil
	}

	// Create session
	svc, err := ShieldClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_shield_attack.getAttack", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	data, err := svc.DescribeAttack(ctx, &shield.DescribeAttackInput{
		AttackId: aws.String(attackId),
	})

	if err != nil {
		plugin.Logger(ctx).Error("aws_shield_attack.getAttack", "api_error", err)
		return nil, err
	}

	if data.Attack != nil {
		return &AttackExtended{
			AttackVectors: getAttackVectors(*data),
			AttackCounters: data.Attack.AttackCounters,
			AttackId: data.Attack.AttackId,
			AttackProperties: data.Attack.AttackProperties,
			EndTime: data.Attack.EndTime,
			Mitigations: data.Attack.Mitigations,
			ResourceArn: data.Attack.ResourceArn,
			StartTime: data.Attack.StartTime,
			SubResources: data.Attack.SubResources,
		}, nil
	}

	return nil, nil
}

func listAttacks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := ShieldClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_shield_attack.listAttacks", "connection_error", err)
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
		input.ResourceArns = []string{}
		for _, q := range d.Quals["resource_arn"].Quals {
			input.ResourceArns = append(input.ResourceArns, q.Value.GetStringValue())
		}
	}

	if d.Quals["start_time"] != nil {
		input.StartTime = getTimeRange(d.Quals["start_time"].Quals)
	}

	if d.Quals["end_time"] != nil {
		input.EndTime = getTimeRange(d.Quals["end_time"].Quals)
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
			plugin.Logger(ctx).Error("aws_shield_attack.listAttacks", "api_error", err)
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