package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsConfigConfigurationRecorder(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_config_configuration_recorder",
		Description: "AWS Config Configuration Recorder",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"NoSuchConfigurationRecorderException"}),
			},
			Hydrate: getConfigConfigurationRecorder,
		},
		List: &plugin.ListConfig{
			Hydrate: listConfigConfigurationRecorders,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
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
	svc, err := ConfigService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &configservice.DescribeConfigurationRecordersInput{}

	// Additonal Filter
	equalQuals := d.KeyColumnQuals
	if equalQuals["name"] != nil {
		input.ConfigurationRecorderNames = []*string{aws.String(equalQuals["name"].GetStringValue())}
	}

	// Pagination not supported as of date
	op, err := svc.DescribeConfigurationRecorders(input)
	if err != nil {
		return nil, err
	}
	if op.ConfigurationRecorders != nil {
		for _, ConfigurationRecorders := range op.ConfigurationRecorders {
			d.StreamListItem(ctx, ConfigurationRecorders)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getConfigConfigurationRecorder(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getConfigConfigurationRecorder")
	quals := d.KeyColumnQuals
	name := quals["name"].GetStringValue()

	// Create Session
	svc, err := ConfigService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &configservice.DescribeConfigurationRecordersInput{
		ConfigurationRecorderNames: []*string{aws.String(name)},
	}
	plugin.Logger(ctx).Trace("paramsparamsparams", "params", params)

	op, err := svc.DescribeConfigurationRecorders(params)
	if err != nil {
		logger.Debug("getConfigConfigurationRecorder", "ERROR", err)
		return nil, err
	}

	if op != nil {
		return op.ConfigurationRecorders[0], nil
	}

	return nil, nil
}

func getConfigConfigurationRecorderStatus(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getConfigConfigurationRecorderStatus")

	configurationRecorder := h.Item.(*configservice.ConfigurationRecorder)

	// Create Session
	svc, err := ConfigService(ctx, d)
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
	region := d.KeyColumnQualString(matrixKeyRegion)

	configurationRecorder := h.Item.(*configservice.ConfigurationRecorder)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	arn := "arn:" + commonColumnData.Partition + ":config:" + region + ":" + commonColumnData.AccountId + ":config-recorder" + "/" + *configurationRecorder.Name

	return arn, nil
}
