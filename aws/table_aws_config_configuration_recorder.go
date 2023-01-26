package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/configservice"
	"github.com/aws/aws-sdk-go-v2/service/configservice/types"

	configservicev1 "github.com/aws/aws-sdk-go/service/configservice"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsConfigConfigurationRecorder(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_config_configuration_recorder",
		Description: "AWS Config Configuration Recorder",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchConfigurationRecorderException"}),
			},
			Hydrate: getConfigConfigurationRecorder,
		},
		List: &plugin.ListConfig{
			Hydrate: listConfigConfigurationRecorders,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(configservicev1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the recorder. By default, AWS Config automatically assigns the name default when creating the configuration recorder.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the configuration recorder.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsConfigurationRecorderARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "recording_group",
				Description: "Specifies the types of AWS resources for which AWS Config records configuration changes.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("RecordingGroup"),
			},
			{
				Name:        "role_arn",
				Description: "Amazon Resource Name (ARN) of the IAM role used to describe the AWS resources associated with the account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RoleARN"),
			},
			{
				Name:        "status_recording",
				Description: "Specifies whether or not the recorder is currently recording.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getConfigConfigurationRecorderStatus,
				Transform:   transform.FromField("Recording"),
			},
			{
				Name:        "status",
				Description: "The current status of the configuration recorder.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getConfigConfigurationRecorderStatus,
				Transform:   transform.FromValue(),
			},
			// Standard columns
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsConfigurationRecorderARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
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

func listConfigConfigurationRecorders(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := ConfigClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_config_configuration_recorder.listConfigConfigurationRecorders", "get_client_error", err)
		return nil, err
	}

	input := &configservice.DescribeConfigurationRecordersInput{}

	// Pagination not supported as of date
	op, err := svc.DescribeConfigurationRecorders(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_config_configuration_recorder.listConfigConfigurationRecorders", "api_error", err)
		return nil, err
	}
	if op.ConfigurationRecorders != nil {
		for _, configurationRecorder := range op.ConfigurationRecorders {
			d.StreamListItem(ctx, configurationRecorder)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getConfigConfigurationRecorder(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	name := quals["name"].GetStringValue()

	// Create session
	svc, err := ConfigClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_config_configuration_recorder.getConfigConfigurationRecorder", "get_client_error", err)
		return nil, err
	}

	params := &configservice.DescribeConfigurationRecordersInput{
		ConfigurationRecorderNames: []string{name},
	}

	op, err := svc.DescribeConfigurationRecorders(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_config_configuration_recorder.getConfigConfigurationRecorder", "api_error", err)
		return nil, err
	}

	if op != nil {
		return op.ConfigurationRecorders[0], nil
	}

	return nil, nil
}

func getConfigConfigurationRecorderStatus(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	configurationRecorder := h.Item.(types.ConfigurationRecorder)

	// Create session
	svc, err := ConfigClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_config_configuration_recorder.getConfigConfigurationRecorderStatus", "get_client_error", err)
		return nil, err
	}

	params := &configservice.DescribeConfigurationRecorderStatusInput{
		ConfigurationRecorderNames: []string{*configurationRecorder.Name},
	}

	status, err := svc.DescribeConfigurationRecorderStatus(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_config_configuration_recorder.getConfigConfigurationRecorderStatus", "api_error", err)
		return nil, err
	}

	if len(status.ConfigurationRecordersStatus) < 1 {
		return nil, nil
	}

	return status.ConfigurationRecordersStatus[0], nil
}

func getAwsConfigurationRecorderARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)

	configurationRecorder := h.Item.(types.ConfigurationRecorder)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_config_configuration_recorder.getAwsConfigurationRecorderARN", "api_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	arn := "arn:" + commonColumnData.Partition + ":config:" + region + ":" + commonColumnData.AccountId + ":config-recorder" + "/" + *configurationRecorder.Name

	return arn, nil
}
