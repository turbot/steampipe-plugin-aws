package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/configservice"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsConfigRetentionConfiguration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_config_retention_configuration",
		Description: "AWS Config Retention Configuration",
		List: &plugin.ListConfig{
			Hydrate: listConfigRetentionConfigurations,
			Tags:    map[string]string{"service": "config", "action": "DescribeRetentionConfigurations"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_CONFIG_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the retention configuration object.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "retention_period_in_days",
				Description: "Number of days Config stores your historical information.",
				Type:        proto.ColumnType_INT,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		}),
	}
}

//// LIST FUNCTION

// Currently, AWS Config supports only one retention configuration per region in your account.
func listConfigRetentionConfigurations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := ConfigClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_config_configuration_recorder.listConfigRetentionConfigurations", "get_client_error", err)
		return nil, err
	}

	input := &configservice.DescribeRetentionConfigurationsInput{}

	paginator := configservice.NewDescribeRetentionConfigurationsPaginator(svc, input, func(o *configservice.DescribeRetentionConfigurationsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_config_configuration_recorder.listConfigRetentionConfigurations", "api_error", err)
			return nil, err
		}

		for _, items := range output.RetentionConfigurations {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

	}

	return nil, nil
}
