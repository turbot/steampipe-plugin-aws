package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/transfer"
	"github.com/aws/aws-sdk-go-v2/service/transfer/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsTransferConnector(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_transfer_connector",
		Description: "AWS Transfer Connector",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("connector_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "ValidationException"}),
			},
			Hydrate: getTransferConnector,
			Tags:    map[string]string{"service": "transfer", "action": "DescribeConnector"},
		},
		List: &plugin.ListConfig{
			Hydrate: listTransferConnectors,
			Tags:    map[string]string{"service": "transfer", "action": "ListConnectors"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getTransferConnector,
				Tags: map[string]string{"service": "transfer", "action": "DescribeConnector"},
			},
			{
				Func: getTransferConnectorTags,
				Tags: map[string]string{"service": "transfer", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_TRANSFER_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "connector_id",
				Description: "The unique identifier for the connector.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the connector.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "url",
				Description: "The URL of the partner's AS2 or SFTP endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "access_role",
				Description: "The Amazon Resource Name (ARN) of the Identity and Access Management role to use.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTransferConnector,
			},
			{
				Name:        "logging_role",
				Description: "The Amazon Resource Name (ARN) of the Identity and Access Management (IAM) role that allows a connector to turn on CloudWatch logging for Amazon S3 events.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTransferConnector,
			},
			{
				Name:        "security_policy_name",
				Description: "The text name of the security policy for the specified connector.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTransferConnector,
			},
			{
				Name:        "service_managed_egress_ip_addresses",
				Description: "The list of egress IP addresses of this connector.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTransferConnector,
			},
			{
				Name:        "as2_config",
				Description: "A structure that contains the parameters for an AS2 connector object.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTransferConnector,
				Transform:   transform.FromField("As2Config"),
			},
			{
				Name:        "sftp_config",
				Description: "A structure that contains the parameters for an SFTP connector object.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTransferConnector,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags that are attached to the connector.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTransferConnectorTags,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTransferConnectorTags,
				Transform:   transform.FromValue().Transform(transferConnectorTagListToTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConnectorId"),
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

func listTransferConnectors(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := TransferClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_transfer_connector.listTransferConnectors", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	maxItems := int32(100)
	params := &transfer.ListConnectorsInput{
		MaxResults: aws.Int32(maxItems),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			params.MaxResults = aws.Int32(limit)
		}
	}

	// Create paginator
	paginator := transfer.NewListConnectorsPaginator(svc, params, func(o *transfer.ListConnectorsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_transfer_connector.listTransferConnectors", "api_error", err)
			return nil, err
		}

		if output != nil && output.Connectors != nil {
			for _, connector := range output.Connectors {
				d.StreamListItem(ctx, connector)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getTransferConnector(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var connectorId string
	if h.Item != nil {
		switch item := h.Item.(type) {
		case types.ListedConnector:
			connectorId = *item.ConnectorId
		case *types.DescribedConnector:
			connectorId = *item.ConnectorId
		}
	} else {
		connectorId = d.EqualsQualString("connector_id")
	}

	if connectorId == "" {
		return nil, nil
	}

	// Create service
	svc, err := TransferClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_transfer_connector.getTransferConnector", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &transfer.DescribeConnectorInput{
		ConnectorId: aws.String(connectorId),
	}

	// Get connector details
	data, err := svc.DescribeConnector(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_transfer_connector.getTransferConnector", "api_error", err)
		return nil, err
	}

	return data.Connector, nil
}

func getTransferConnectorTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var connectorArn string
	if h.Item != nil {
		switch item := h.Item.(type) {
		case types.ListedConnector:
			connectorArn = *item.Arn
		case *types.DescribedConnector:
			connectorArn = *item.Arn
		}
	} else {
		// For get call, we need to get the ARN from the connector details
		connector, err := getTransferConnector(ctx, d, h)
		if err != nil {
			return nil, err
		}
		if connector == nil {
			return nil, nil
		}
		connectorArn = *connector.(*types.DescribedConnector).Arn
	}

	if connectorArn == "" {
		return nil, nil
	}

	// Create service
	svc, err := TransferClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_transfer_connector.getTransferConnectorTags", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &transfer.ListTagsForResourceInput{
		Arn: aws.String(connectorArn),
	}

	// Get connector tags
	data, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_transfer_connector.getTransferConnectorTags", "api_error", err)
		return nil, err
	}

	return data, nil
}

//// TRANSFORM FUNCTIONS

func transferConnectorTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tagList := d.HydrateItem.(*transfer.ListTagsForResourceOutput)

	if len(tagList.Tags) > 0 {
		turbotTagsMap := map[string]string{}
		for _, i := range tagList.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
