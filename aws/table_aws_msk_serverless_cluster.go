package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/kafka/types"

	kafkav1 "github.com/aws/aws-sdk-go/service/kafka"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsMSKServerlessCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_msk_serverless_cluster",
		Description: "AWS Serverless Managed Streaming for Apache Kafka",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			Hydrate:    getKafkaCluster(string(types.ClusterTypeServerless)),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listKafkaClusters(string(types.ClusterTypeServerless)),
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
				Name:        "cluster_operation",
				Description: "Description of this MSK operation.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKafkaClusterOperation,
				Transform:   transform.FromValue(),
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
