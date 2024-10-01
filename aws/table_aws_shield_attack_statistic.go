package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/shield"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAwsShieldAttackStatistic(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_shield_attack_statistic",
		Description: "AWS Shield Attack Statistic",
		List: &plugin.ListConfig{
			Hydrate: listAwsShieldAttackStatistic,
			Tags:    map[string]string{"service": "shield", "action": "DescribeAttackStatistics"},
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "start_time",
				Description: "The start time of observation time range (should be always one year ago).",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "end_time",
				Description: "The end time of the observation time range (should be always the current date).",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "unit",
				Description: "Unit of the attack statistic.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "max",
				Description: "The maximum attack volume observed in the observation time range for the given unit.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "attack_count",
				Description: "The number of attacks detected during the time period. This is always present, but might be zero.",
				Type:        proto.ColumnType_INT,
			},
		}),
	}
}

type attackStatistic struct {
	StartTime time.Time
    EndTime time.Time
	Unit string
	Max float64
	AttackCount int64
}

//// HYDRATE FUNCTIONS

func listAwsShieldAttackStatistic(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := ShieldClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_shield_attack_statistic.listAwsShieldAttackStatistic", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	data, err := svc.DescribeAttackStatistics(ctx, &shield.DescribeAttackStatisticsInput{})

	if err != nil {
		plugin.Logger(ctx).Error("aws_shield_attack_statistic.getAwsShieldAttackStatistic", "api_error", err)
		return nil, err
	}

	for _, stat := range data.DataItems {
		var unit string
		var max float64

		if stat.AttackCount == 0 {
			continue
		} else if stat.AttackVolume.BitsPerSecond != nil {
			unit = "BitsPerSecond"
			max = stat.AttackVolume.BitsPerSecond.Max
		} else if stat.AttackVolume.PacketsPerSecond != nil {
			unit = "PacketsPerSecond"
			max = stat.AttackVolume.PacketsPerSecond.Max
		} else if stat.AttackVolume.RequestsPerSecond != nil {
			unit = "RequestsPerSecond"
			max = stat.AttackVolume.RequestsPerSecond.Max
		}

		d.StreamListItem(ctx, &attackStatistic{
			StartTime:   *data.TimeRange.FromInclusive,
			EndTime:   	 *data.TimeRange.ToExclusive,
			Unit:        unit,
			Max:         max,
			AttackCount: stat.AttackCount,
		})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}
