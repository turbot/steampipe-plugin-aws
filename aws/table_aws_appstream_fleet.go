package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/appstream"
	"github.com/aws/aws-sdk-go-v2/service/appstream/types"
	appstreamEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"
	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAppStreamFleet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_appstream_fleet",
		Description: "AWS AppStream Fleet",
		List: &plugin.ListConfig{
			Hydrate: listAppStreamFleets,
			Tags:    map[string]string{"service": "appstream", "action": "DescribeFleets"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAppStreamFleetTags,
				Tags: map[string]string{"service": "appstream", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(appstreamEndpoint.APPSTREAM2ServiceID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the fleet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the fleet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_type",
				Description: "The instance type to use when launching fleet instances.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The current state for the fleet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_time",
				Description: "The time the fleet was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "The description to display.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "The fleet name to display.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "disconnect_timeout_in_seconds",
				Description: "The amount of time that a streaming session remains active after users disconnect. If they try to reconnect to the streaming session after a disconnection or network interruption within this time interval, they are connected to their previous session. Otherwise, they are connected to a new session with a new streaming instance. Specify a value between 60 and 360000.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "directory_name",
				Description: "The fully qualified name of the directory (for example, corp.example.com).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DomainJoinInfo.DirectoryName"),
			},
			{
				Name:        "organizational_unit_distinguished_name",
				Description: "The distinguished name of the organizational unit for computer accounts.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DomainJoinInfo.OrganizationalUnitDistinguishedName"),
			},
			{
				Name:        "enable_default_internet_access",
				Description: "Indicates whether default internet access is enabled for the fleet.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "fleet_type",
				Description: "The fleet type. ALWAYS_ON Provides users with instant-on access to their apps. You are charged for all running instances in your fleet, even if no users are streaming apps. ON_DEMAND Provide users with access to applications after they connect, which takes one to two minutes. You are charged for instance streaming when users are connected and a small hourly fee for instances that are not streaming apps.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "iam_role_arn",
				Description: "The ARN of the IAM role that is applied to the fleet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "idle_disconnect_timeout_in_seconds",
				Description: "The amount of time that users can be idle (inactive) before they are disconnected from their streaming session and the DisconnectTimeoutInSeconds time interval begins.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "image_arn",
				Description: "The ARN for the public, private, or shared image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "image_name",
				Description: "The name of the image used to create the fleet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "max_concurrent_sessions",
				Description: "The maximum number of concurrent sessions for the fleet.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "max_user_duration_in_seconds",
				Description: "The maximum amount of time that a streaming session can remain active, in seconds. If users are still connected to a streaming instance five minutes before this limit is reached, they are prompted to save any open documents before being disconnected.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "platform",
				Description: "The platform of the fleet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "stream_view",
				Description: "The AppStream 2.0 view that is displayed to your users when they stream from the fleet. When APP is specified, only the windows of applications opened by users display. When DESKTOP is specified, the standard desktop that is provided by the operating system displays. The default value is APP.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "compute_capacity_status",
				Description: "The capacity status for the fleet.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "fleet_errors",
				Description: "The fleet errors.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "session_script_s3_location",
				Description: "The S3 location of the session scripts configuration zip file. This only applies to Elastic fleets.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "usb_device_filter_strings",
				Description: "The USB device filter strings associated with the fleet.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "vpc_config",
				Description: "The VPC configuration for the fleet.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAppStreamFleetTags,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listAppStreamFleets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create Session
	svc, err := AppStreamClient(ctx, d)
	if err != nil {
		logger.Error("aws_appstream_fleet.listAppStreamFleets", "connection_error", err)
		return nil, err
	}

	// Unsupported region check
	if svc == nil {
		return nil, nil
	}

	params := &appstream.DescribeFleetsInput{}

	if d.Quals["name"] != nil {
		for _, q := range d.Quals["name"].Quals {
			value := q.Value.GetStringValue()
			if q.Operator == "=" {
				params.Names = append(params.Names, value)
			}
		}
	}

	pageLeft := true

	for pageLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		op, err := svc.DescribeFleets(ctx, params)

		if err != nil {
			logger.Error("aws_appstream_fleet.listAppStreamFleets", "api_error", err)
			return nil, err
		}

		for _, fleet := range op.Fleets {
			d.StreamListItem(ctx, fleet)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if op.NextToken != nil {
			params.NextToken = op.NextToken
		} else {
			break
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAppStreamFleetTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	var arn string
	if h.Item != nil {
		arn = *h.Item.(types.Fleet).Arn
	} else {
		return nil, nil
	}

	// Create Session
	svc, err := AppStreamClient(ctx, d)
	if err != nil {
		logger.Error("aws_appstream_fleet.getAppStreamFleetTags", "connection_error", err)
		return nil, err
	}

	params := &appstream.ListTagsForResourceInput{
		ResourceArn: &arn,
	}

	tags, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "ResourceNotFoundException" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_appstream_fleet.getAppStreamFleetTags", "api_error", err)
		return nil, err
	}

	return tags.Tags, nil
}
