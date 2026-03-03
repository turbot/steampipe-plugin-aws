package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/servicequotas"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsServiceQuotasAutoManagementConfiguration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_servicequotas_auto_management_configuration",
		Description: "AWS Service Quotas Auto Management Configuration",
		DefaultIgnoreConfig: &plugin.IgnoreConfig{
			ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchResourceException"}),
		},
		List: &plugin.ListConfig{
			Hydrate: getAutoManagementConfiguration,
			Tags:    map[string]string{"service": "servicequotas", "action": "GetAutoManagementConfiguration"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_SERVICEQUOTAS_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "opt_in_status",
				Description: "Status on whether Automatic Management is started or stopped.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "opt_in_type",
				Description: "Information on the opt-in type for Automatic Management. There are two modes: `NotifyOnly` (notify only) and `NotifyAndAdjust` (notify and auto-adjust).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "opt_in_level",
				Description: "Information on the opt-in level for Automatic Management. Only ACCOUNT level is supported.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "notification_arn",
				Description: "The User Notifications Amazon Resource Name (ARN) for Automatic Management notifications.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "exclusion_list",
				Description: "List of Amazon Web Services services excluded from Automatic Management.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromConstant("Auto Management Configuration"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAutoManagementConfigurationAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func getAutoManagementConfiguration(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := ServiceQuotasClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_servicequotas_auto_management_configuration.getAutoManagementConfiguration", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	result, err := svc.GetAutoManagementConfiguration(ctx, &servicequotas.GetAutoManagementConfigurationInput{})
	if err != nil {
		plugin.Logger(ctx).Error("aws_servicequotas_auto_management_configuration.getAutoManagementConfiguration", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, result)
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAutoManagementConfigurationAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_servicequotas_auto_management_configuration.getAutoManagementConfigurationAkas", "common_data_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	arn := fmt.Sprintf("arn:%s:servicequotas:%s:%s:auto-management-configuration", commonColumnData.Partition, region, commonColumnData.AccountId)

	return []string{arn}, nil
}
