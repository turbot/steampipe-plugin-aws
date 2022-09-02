package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kafka"
	"github.com/aws/aws-sdk-go-v2/service/kafka/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsMSKCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_msk_cluster",
		Description: "AWS Managed Streaming for Apache Kafka",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("cluster_arn"),
			// IgnoreConfig: &plugin.IgnoreConfig{
			// 	ShouldIgnoreErrorFunc: isNotFoundErrorV2([]string{"DBClusterNotFoundFault"}),
			// },
			Hydrate: getKafkaCluster,
		},
		List: &plugin.ListConfig{
			Hydrate: listKafkaClusters,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "cluster_arn",
				Description: "The Amazon Resource Name (ARN) that uniquely identifies the Cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_name",
				Description: "The name of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "active_operation_arn",
				Description: "Arn of active cluster operation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_type",
				Description: "The type of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The time when the cluster was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "current_version",
				Description: "The current version of the MSK cluster.",
				Type:        proto.ColumnType_STRING,
			},
			// {
			// 	Name:        "number_of_broker_nodes",
			// 	Description: "The number of broker nodes in the cluster.",
			// 	Type:        proto.ColumnType_TIMESTAMP,
			// },
			{
				Name:        "state",
				Description: "Settings for open monitoring using Prometheus.",
				Type:        proto.ColumnType_STRING,
			},
			// {
			// 	Name:        "zookeeper_connect_string",
			// 	Description: "The connection string to use to connect to the Apache ZooKeeper cluster.",
			// 	Type:        proto.ColumnType_TIMESTAMP,
			// },
			// {
			// 	Name:        "zookeeper_connect_string_tls",
			// 	Description: "The connection string to use to connect to zookeeper cluster on Tls port.",
			// 	Type:        proto.ColumnType_TIMESTAMP,
			// },

			// JSON columns
			// {
			// 	Name:        "broker_node_group_info",
			// 	Description: "Information about the broker nodes.",
			// 	Type:        proto.ColumnType_JSON,
			// },
			// {
			// 	Name:        "client_authentication",
			// 	Description: "Includes all client authentication information.",
			// 	Type:        proto.ColumnType_JSON,
			// },
			// {
			// 	Name:        "current_broker_software_info",
			// 	Description: "Information about the version of software currently deployed on the Apache Kafka brokers in the cluster.",
			// 	Type:        proto.ColumnType_JSON,
			// },
			// {
			// 	Name:        "encryption_info",
			// 	Description: "Includes all encryption-related information.",
			// 	Type:        proto.ColumnType_JSON,
			// },
			// {
			// 	Name:        "enhanced_monitoring",
			// 	Description: "Specifies which metrics are gathered for the MSK cluster.",
			// 	Type:        proto.ColumnType_JSON,
			// },
			// {
			// 	Name:        "logging_info",
			// 	Description: "Includes all logging-related information.",
			// 	Type:        proto.ColumnType_JSON,
			// },
			// {
			// 	Name:        "open_monitoring",
			// 	Description: "Settings for open monitoring using Prometheus.",
			// 	Type:        proto.ColumnType_JSON,
			// },
			// {
			// 	Name:        "cluster_configuration",
			// 	Description: "Description of this MSK configuration.",
			// 	Type:        proto.ColumnType_JSON,
			// 	Hydrate:     getKafkaClusterConfiguration,
			// 	Transform:   transform.FromValue(),
			// },
			{
				Name:        "cluster_operation",
				Description: "Description of this MSK operation.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKafkaClusterOperation,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "provisioned",
				Description: "Information about the provisioned cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "serverless",
				Description: "Information about the serverless cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "state_info",
				Description: "State Info for the Amazon MSK cluster.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard columns
			{
				Name:        "tags",
				Description: "A list of tags attached to the Cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DBClusterIdentifier"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DBClusterArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listKafkaClusters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create Session
	svc, err := KafkaClient(ctx, d)
	if err != nil {
		logger.Error("aws_msk_cluster.listKafkaClusters", "service_creation_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 20 {
				maxLimit = 20
			} else {
				maxLimit = limit
			}
		}
	}

	input := kafka.ListClustersV2Input{
		MaxResults: maxLimit,
	}

	paginator := kafka.NewListClustersV2Paginator(svc, &input, func(o *kafka.ListClustersV2PaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_msk_cluster.listKafkaClusters", "api_error", err)
			return nil, err
		}

		for _, cluster := range output.ClusterInfoList {
			d.StreamListItem(ctx, cluster)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getKafkaCluster(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	clusterArn := d.KeyColumnQuals["cluster_arn"].GetStringValue()
	if len(clusterArn) < 1 {
		return nil, nil
	}

	// Create service
	svc, err := KafkaClient(ctx, d)
	if err != nil {
		logger.Error("aws_msk_cluster.getKafkaCluster", "service_creation_error", err)
		return nil, err
	}

	params := &kafka.DescribeClusterV2Input{
		ClusterArn: aws.String(clusterArn),
	}

	op, err := svc.DescribeClusterV2(ctx, params)
	if err != nil {
		logger.Error("aws_msk_cluster.getKafkaCluster", "api_error", err)
		return nil, err
	}

	return op.ClusterInfo, nil
}

func getKafkaClusterConfiguration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	cluster := h.Item.(types.Cluster)

	// Create Session
	svc, err := KafkaClient(ctx, d)
	if err != nil {
		logger.Error("aws_msk_cluster.getKafkaClusterConfiguration", "service_creation_error", err)
		return nil, err
	}

	// Build the params
	params := &kafka.DescribeConfigurationInput{
		Arn: cluster.ClusterArn,
	}

	// Get call
	op, err := svc.DescribeConfiguration(ctx, params)
	if err != nil {
		logger.Error("aws_kafka_cluster.getKafkaClusterConfiguration", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getKafkaClusterOperation(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	cluster := h.Item.(types.Cluster)

	if cluster.ActiveOperationArn == nil {
		return nil, nil
	}

	// Create Session
	svc, err := KafkaClient(ctx, d)
	if err != nil {
		logger.Error("aws_msk_cluster.getKafkaClusterOperation", "service_creation_error", err)
		return nil, err
	}

	// Build the params
	params := &kafka.DescribeClusterOperationInput{
		ClusterOperationArn: cluster.ActiveOperationArn,
	}

	// Get call
	op, err := svc.DescribeClusterOperation(ctx, params)
	if err != nil {
		logger.Error("aws_kafka_cluster.getKafkaClusterOperation", "api_error", err)
		return nil, err
	}

	return op.ClusterOperationInfo, nil
}
