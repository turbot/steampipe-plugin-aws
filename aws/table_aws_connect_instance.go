package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/connect"
	"github.com/aws/aws-sdk-go-v2/service/connect/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsConnectInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_connect_instance",
		Description: "AWS Connect Instance",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getConnectInstance,
			Tags:    map[string]string{"service": "connect", "action": "DescribeInstance"},
		},
		List: &plugin.ListConfig{
			Hydrate: listConnectInstances,
			Tags:    map[string]string{"service": "connect", "action": "ListInstances"},
		},

		GetMatrixItemFunc: SupportedRegionMatrix(AWS_CONNECT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The identifier of the Amazon Connect instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_alias",
				Description: "The alias of the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_status",
				Description: "The state of the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "identity_management_type",
				Description: "The identity management type of the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_access_url",
				Description: "This URL allows contact center users to access the Amazon Connect admin website.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_role",
				Description: "The service role of the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_time",
				Description: "When the instance was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "status_reason",
				Description: "Relevant details why the instance was not successfully created.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getConnectInstance,
			},
			{
				Name:        "inbound_calls_enabled",
				Description: "Whether inbound calls are enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "outbound_calls_enabled",
				Description: "Whether outbound calls are enabled.",
				Type:        proto.ColumnType_BOOL,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InstanceAlias"),
			},
			{
				Name:        "tags",
				Description: "The tags of the instance.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getConnectInstance,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getConnectInstance,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listConnectInstances(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := ConnectClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_connect_instance.listConnectInstances", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	maxResults := int32(1000)
	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil && *d.QueryContext.Limit < int64(maxResults) {
		maxResults = int32(*d.QueryContext.Limit)
	}

	input := &connect.ListInstancesInput{
		MaxResults: &maxResults,
	}

	paginator := connect.NewListInstancesPaginator(svc, input, func(o *connect.ListInstancesPaginatorOptions) {
		o.Limit = maxResults
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_connect_instance.listConnectInstances", "api_error", err)
			return nil, err
		}

		for _, instance := range output.InstanceSummaryList {
			d.StreamListItem(ctx, instance)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getConnectInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var instanceId string
	if h.Item != nil {
		// If we have an item from the list, extract the instance ID
		if instance, ok := h.Item.(types.InstanceSummary); ok {
			instanceId = *instance.Id
		}
	} else {
		// If this is a get call, use the key column
		instanceId = d.EqualsQualString("id")
	}

	if instanceId == "" {
		return nil, nil
	}

	// Create service
	svc, err := ConnectClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_connect_instance.getConnectInstance", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &connect.DescribeInstanceInput{
		InstanceId: aws.String(instanceId),
	}

	// Get call
	data, err := svc.DescribeInstance(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_connect_instance.getConnectInstance", "api_error", err)
		return nil, err
	}

	if data != nil && data.Instance != nil {
		return data.Instance, nil
	}

	return nil, nil
}
