package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSSMServiceSetting(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssm_service_setting",
		Description: "AWS SSM Service Setting",
		List: &plugin.ListConfig{
			KeyColumns: plugin.AllColumns([]string{"setting_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationException", "InvalidParameterValue"}),
			},
			Hydrate: getSSMServiceSetting,
			Tags:    map[string]string{"service": "ssm", "action": "GetServiceSetting"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_SSM_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "setting_id",
				Description: "The ID of the service setting.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The ARN of the service setting.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ARN"),
			},
			{
				Name:        "setting_value",
				Description: "The value of the service setting.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the service setting. The value can be Default, Customized or PendingUpdate.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_modified_date",
				Description: "The last time the service setting was modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_modified_user",
				Description: "The ARN of the last modified user. This field is populated only if the setting value was overwritten.",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe Standard Columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SettingId"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ARN").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// GET FUNCTION

func getSSMServiceSetting(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	settingID := d.EqualsQuals["setting_id"].GetStringValue()

	if settingID == "" {
		return nil, nil
	}

	// Create service
	svc, err := SSMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_service_setting.getSSMServiceSetting", "connection_error", err)
		return nil, err
	}

	params := &ssm.GetServiceSettingInput{
		SettingId: aws.String(settingID),
	}

	d.WaitForListRateLimit(ctx)
	op, err := svc.GetServiceSetting(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_service_setting.getSSMServiceSetting", "api_error", err)
		return nil, err
	}

	if op.ServiceSetting != nil {
		d.StreamListItem(ctx, op.ServiceSetting)
	}

	return nil, nil
}
