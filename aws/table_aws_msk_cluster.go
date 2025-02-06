package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kafka"
	"github.com/aws/aws-sdk-go-v2/service/kafka/types"

	kafkav1 "github.com/aws/aws-sdk-go/service/kafka"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsMSKCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_msk_cluster",
		Description: "AWS Managed Streaming for Apache Kafka",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			Hydrate:    getKafkaCluster(string(types.ClusterTypeProvisioned)),
			Tags:       map[string]string{"service": "kafka", "action": "DescribeCluster"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listKafkaClusters(string(types.ClusterTypeProvisioned)),
			Tags:    map[string]string{"service": "kafka", "action": "ListClusters"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getKafkaClusterConfiguration,
				Tags: map[string]string{"service": "kafka", "action": "DescribeConfiguration"},
			},
			{
				Func: getKafkaClusterOperation,
				Tags: map[string]string{"service": "kafka", "action": "DescribeClusterOperation"},
			},
			{
				Func: getKafkaClusterBootstrapBrokers,
				Tags: map[string]string{"service": "kafka", "action": "GetBootstrapBrokers"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(kafkav1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) that uniquely identifies the Cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterArn"),
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
			{
				Name:        "state",
				Description: "Settings for open monitoring using Prometheus.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_configuration",
				Description: "Description of this MSK configuration.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKafkaClusterConfiguration,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "cluster_operation",
				Description: "Description of this MSK operation.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKafkaClusterOperation,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "bootstrap_broker_string",
				Description: "A string containing one or more hostname:port pairs of Kafka brokers suitable for use with Apache Kafka clients.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKafkaClusterBootstrapBrokers,
			},
			{
				Name:        "bootstrap_broker_tls",
				Description: "A string containing one or more hostname:port pairs of Kafka brokers suitable for TLS authentication.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKafkaClusterBootstrapBrokers,
			},
			{
				Name:        "provisioned",
				Description: "Information about the provisioned cluster.",
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
				Transform:   transform.FromField("ClusterName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listKafkaClusters(clusterType string) func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		logger := plugin.Logger(ctx)

		// Create Session
		svc, err := KafkaClient(ctx, d)
		if err != nil {
			logger.Error("aws_msk_cluster.listKafkaClusters", "service_creation_error", err)
			return nil, err
		}
		if svc == nil {
			// Unsupported region, return no data
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
			MaxResults:        aws.Int32(maxLimit),
			ClusterTypeFilter: &clusterType,
		}

		paginator := kafka.NewListClustersV2Paginator(svc, &input, func(o *kafka.ListClustersV2PaginatorOptions) {
			o.Limit = maxLimit
			o.StopOnDuplicateToken = true
		})

		for paginator.HasMorePages() {
			// apply rate limiting
			d.WaitForListRateLimit(ctx)

			output, err := paginator.NextPage(ctx)
			if err != nil {
				plugin.Logger(ctx).Error("aws_msk_cluster.listKafkaClusters", "api_error", err)
				return nil, err
			}

			for _, cluster := range output.ClusterInfoList {
				d.StreamListItem(ctx, cluster)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}

		return nil, nil
	}
}

//// HYDRATE FUNCTIONS

func getKafkaCluster(clusterType string) func(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		logger := plugin.Logger(ctx)
		clusterArn := d.EqualsQuals["arn"].GetStringValue()
		if clusterArn == "" {
			return nil, nil
		}

		// Create service
		svc, err := KafkaClient(ctx, d)
		if err != nil {
			logger.Error("aws_msk_cluster.getKafkaCluster", "service_creation_error", err)
			return nil, err
		}

		// Unsupported region, return no data
		if svc == nil {
			return nil, nil
		}

		params := &kafka.DescribeClusterV2Input{
			ClusterArn: aws.String(clusterArn),
		}

		op, err := svc.DescribeClusterV2(ctx, params)
		if err != nil {
			logger.Error("aws_msk_cluster.getKafkaCluster", "api_error", err)
			return nil, err
		}

		// It'd be better to check if the cluster type matches before the API call,
		// but we can't tell what type of cluster it is based off of the ARN.
		if op != nil && string(op.ClusterInfo.ClusterType) == clusterType {
			return *op.ClusterInfo, nil
		}
		return nil, nil
	}
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
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
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

func getKafkaClusterConfiguration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	cluster := h.Item.(types.Cluster)
	var configArn string
	if cluster.ClusterType == types.ClusterTypeProvisioned {
		if cluster.Provisioned.CurrentBrokerSoftwareInfo.ConfigurationArn != nil {
			configArn = *cluster.Provisioned.CurrentBrokerSoftwareInfo.ConfigurationArn
		}
	}

	if len(configArn) < 1 {
		return nil, nil
	}

	// Create Session
	svc, err := KafkaClient(ctx, d)
	if err != nil {
		logger.Error("aws_msk_cluster.getKafkaClusterConfiguration", "service_creation_error", err)
		return nil, err
	}
	// Unsupported region, return no data
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &kafka.DescribeConfigurationInput{
		Arn: &configArn,
	}

	// Get call
	op, err := svc.DescribeConfiguration(ctx, params)
	if err != nil {
		logger.Error("aws_kafka_cluster.getKafkaClusterConfiguration", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getKafkaClusterBootstrapBrokers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	cluster := h.Item.(types.Cluster)

	clusterArn := aws.ToString(cluster.ClusterArn)
	// Empty check
	if clusterArn == "" {
		return nil, nil
	}

	svc, err := KafkaClient(ctx, d)
	if err != nil {
		logger.Error("aws_msk_cluster.getKafkaClusterBootstrapBrokers", "service_creation_error", err)
		return nil, err
	}

	// Unsupported region, return no data
	if svc == nil {
		return nil, nil
	}

	params := &kafka.GetBootstrapBrokersInput{
		ClusterArn: &clusterArn,
	}

	op, err := svc.GetBootstrapBrokers(ctx, params)
	if err != nil {
		logger.Error("aws_msk_cluster.getKafkaClusterBootstrapBrokers", "api_error", err)
		return nil, err
	}

	return op, nil
}
