package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/shield"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsShieldAttack(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_shield_attack",
		Description: "AWS Shield Attack",
		List: &plugin.ListConfig{
			Hydrate: listAwsShieldAttack,
			KeyColumns: plugin.SingleColumn("attack_id"),
			Tags:    map[string]string{"service": "shield", "action": "DescribeAttack"},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "attack_id",
				Description: "The unique identifier (ID) for the attack.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("attack_id"),
			},
			{
				Name:        "resource_arn",
				Description: "The ARN (Amazon Resource Name) of the resource that was attacked.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceArn"),
			},
			{
				Name:        "sub_resources",
				Description: "If applicable, additional detail about the resource being attacked, for example, IP address or URL.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SubResources"),
			},
			{
				Name:        "start_time",
				Description: "The time the attack started.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("StartTime"),
			},
			{
				Name:        "end_time",
				Description: "The time the attack ended.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("EndTime"),
			},
			{
				Name:        "attack_counters",
				Description: "List of counters that describe the attack for the specified time period.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AttackCounters"),
			},
			{
				Name:        "attack_properties",
				Description: "The array of objects that provide details of the Shield event.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AttackProperties"),
			},
			{
				Name:        "mitigations",
				Description: "List of mitigation actions taken for the attack.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Mitigations"),
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

func listAwsShieldAttack(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := ShieldClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_shield_attack.listAwsShieldAttack", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	input := &shield.DescribeAttackInput{
		AttackId: aws.String(d.EqualsQualString("attack_id")),
	}

	data, err := svc.DescribeAttack(ctx, input)

	if err != nil {
		plugin.Logger(ctx).Error("aws_shield_attack.getAwsShieldAttack", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, data.Attack)

	return nil, nil
}
