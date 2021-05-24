package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/go-kit/types"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsConfigConfigurationRecorder(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_config_configuration_recorder",
		Description: "AWS Config Configuration Recorder",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "region"}),
			ShouldIgnoreError: isNotFoundError([]string{"NoSuchConfigurationRecorderException", "UnrecognizedClientException"}),
			Hydrate:           getConfigConfigurationRecorder,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAwsRegions,
			Hydrate: listConfigConfigurationRecorders,
		},
		Columns: awsS3Columns([]*plugin.Column{
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
			{
				Name:        "region",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_STRING,
			},
		}),
	}
}

type recorderInfo = struct {
	configservice.ConfigurationRecorder
	Region string
}

//// LIST FUNCTION

func listConfigConfigurationRecorders(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := h.Item.(*ec2.Region)

	// If a region is not opted-in, we cannot list the availability zones
	if types.SafeString(region.OptInStatus) == "not-opted-in" {
		return nil, nil
	}

	plugin.Logger(ctx).Trace("listConfigConfigurationRecorders", "AWS_REGION", region)

	// Create session
	svc, err := ConfigService(ctx, d, *region.RegionName)
	if err != nil {
		return nil, err
	}

	op, err := svc.DescribeConfigurationRecorders(
		&configservice.DescribeConfigurationRecordersInput{})
	if err != nil {
		return nil, err
	}
	if op.ConfigurationRecorders != nil {
		for _, ConfigurationRecorders := range op.ConfigurationRecorders {
			d.StreamLeafListItem(ctx, recorderInfo{*ConfigurationRecorders, *region.RegionName })
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getConfigConfigurationRecorder(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getConfigConfigurationRecorder")
	quals := d.KeyColumnQuals
	name := quals["name"].GetStringValue()
	region := quals["region"].GetStringValue()

	plugin.Logger(ctx).Trace("AWS_REGION", "Region", region)

	// Create Session
	svc, err := ConfigService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &configservice.DescribeConfigurationRecordersInput{
		ConfigurationRecorderNames: []*string{aws.String(name)},
	}

	op, err := svc.DescribeConfigurationRecorders(params)
	if err != nil {
		return nil, err
	}

	if op != nil {
		return recorderInfo{*op.ConfigurationRecorders[0], region }, nil
	}
	return nil, nil
}

func getConfigConfigurationRecorderStatus(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getConfigConfigurationRecorderStatus")
	region := h.Item.(recorderInfo).Region
	plugin.Logger(ctx).Trace("AWS_REGION", "Region", region)

	configurationRecorder := h.Item.(recorderInfo).ConfigurationRecorder

	// Create Session
	svc, err := ConfigService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &configservice.DescribeConfigurationRecorderStatusInput{
		ConfigurationRecorderNames: []*string{configurationRecorder.Name},
	}

	status, err := svc.DescribeConfigurationRecorderStatus(params)
	if err != nil {
		return nil, err
	}

	return status.ConfigurationRecordersStatus[0], nil
}

//// TRANSFORM FUNCTIONS

func getAwsConfigurationRecorderARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsConfigurationRecorderAkas")
	configurationRecorder := h.Item.(recorderInfo).ConfigurationRecorder
	regionName := h.Item.(recorderInfo).Region
	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	arn := "arn:" + commonColumnData.Partition + ":config:" + regionName + ":" + commonColumnData.AccountId + ":config-recorder" + "/" + *configurationRecorder.Name

	return arn, nil
}
