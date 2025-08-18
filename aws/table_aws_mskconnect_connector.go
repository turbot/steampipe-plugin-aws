package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kafkaconnect"
	"github.com/aws/aws-sdk-go-v2/service/kafkaconnect/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsMSKConnectConnector(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_mskconnect_connector",
		Description: "AWS MSK Connect Connector",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException", "BadRequestException"}),
			},
			Hydrate: getMSKConnectConnector,
			Tags:    map[string]string{"service": "kafkaconnect", "action": "DescribeConnector"},
		},
		List: &plugin.ListConfig{
			Hydrate: listMSKConnectConnectors,
			Tags:    map[string]string{"service": "kafkaconnect", "action": "ListConnectors"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getMSKConnectConnector,
				Tags: map[string]string{"service": "kafkaconnect", "action": "DescribeConnector"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_KAFKACONNECT_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "connector_name",
				Description: "The name of the connector.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the connector.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMSKConnectConnector,
				Transform:   transform.FromField("ConnectorArn"),
			},
			{
				Name:        "connector_description",
				Description: "The description of the connector.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "connector_state",
				Description: "The state of the connector.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The time that the connector was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "current_version",
				Description: "The current version of the connector.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kafka_connect_version",
				Description: "The version of Kafka Connect. It has to be compatible with both the Apache Kafka cluster's version and the plugins.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_execution_role_arn",
				Description: "The Amazon Resource Name (ARN) of the IAM role used by the connector to access Amazon Web Services resources.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "capacity",
				Description: "The connector's compute capacity settings.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "kafka_cluster",
				Description: "The details of the Apache Kafka cluster to which the connector is connected.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "kafka_cluster_client_authentication",
				Description: "The type of client authentication used to connect to the Apache Kafka cluster. The value is NONE when no client authentication is used.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "kafka_cluster_encryption_in_transit",
				Description: "Details of encryption in transit to the Apache Kafka cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "log_delivery",
				Description: "The settings for delivering connector logs to Amazon CloudWatch Logs.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "plugins",
				Description: "Specifies which plugins were used for this connector.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "worker_configuration",
				Description: "The worker configurations that are in use with the connector.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "connector_configuration",
				Description: "A map of keys to values that represent the configuration for the connector.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMSKConnectConnector,
			},
			{
				Name:        "state_description",
				Description: "Details about the state of a connector.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMSKConnectConnector,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ConnectorName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ConnectorArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listMSKConnectConnectors(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := KafkaConnectClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_mskconnect_connector.listMSKConnectConnectors", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	maxItems := int32(100)
	params := &kafkaconnect.ListConnectorsInput{
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
	paginator := kafkaconnect.NewListConnectorsPaginator(svc, params, func(o *kafkaconnect.ListConnectorsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_mskconnect_connector.listMSKConnectConnectors", "api_error", err)
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

func getMSKConnectConnector(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var connectorArn string
	if h.Item != nil {
		connectorArn = *h.Item.(types.ConnectorSummary).ConnectorArn
	} else {
		connectorArn = d.EqualsQualString("arn")
	}

	if connectorArn == "" {
		return nil, nil
	}

	// Create service
	svc, err := KafkaConnectClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_mskconnect_connector.getMSKConnectConnector", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &kafkaconnect.DescribeConnectorInput{
		ConnectorArn: aws.String(connectorArn),
	}

	// Get connector details
	data, err := svc.DescribeConnector(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_mskconnect_connector.getMSKConnectConnector", "api_error", err)
		return nil, err
	}

	return data, nil
}
