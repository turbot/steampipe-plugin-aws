package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEksCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_eks_cluster",
		Description: "AWS Elastic Kubernetes Service Cluster",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getEksCluster,
		},
		List: &plugin.ListConfig{
			Hydrate: listEksClusters,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the cluster.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksCluster,
			},
			{
				Name:        "created_at",
				Description: "The Unix epoch timestamp in seconds for when the cluster was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getEksCluster,
			},
			{
				Name:        "version",
				Description: "The Kubernetes server version for the cluster.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksCluster,
			},
			{
				Name:        "endpoint",
				Description: "The endpoint for your Kubernetes API server.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksCluster,
			},
			{
				Name:        "role_arn",
				Description: "The Amazon Resource Name (ARN) of the IAM role that provides permissions for the Kubernetes control plane to make calls to AWS API operations on your behalf.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksCluster,
			},
			{
				Name:        "encryption_config",
				Description: "The encryption configuration for the cluster.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksCluster,
			},
			{
				Name:        "resources_vpc_config",
				Description: "The VPC configuration used by the cluster control plane.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksCluster,
			},
			{
				Name:        "kubernetes_network_config",
				Description: "The Kubernetes network configuration for the cluster.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksCluster,
			},
			{
				Name:        "logging",
				Description: "The logging configuration for the cluster.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksCluster,
			},
			{
				Name:        "identity",
				Description: "The identity provider information for the cluster.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksCluster,
			},
			{
				Name:        "status",
				Description: "The current status of the cluster.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksCluster,
			},
			{
				Name:        "certificate_authority",
				Description: "The certificate-authority-data for the cluster.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksCluster,
			},
			{
				Name:        "platform_version",
				Description: "The platform version of your Amazon EKS cluster.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksCluster,
			},
			{
				Name:        "tags",
				Description: "A list of tags assigned to the table",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksCluster,
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
				Hydrate:     getEksCluster,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(arnToAkas),
				Hydrate:     getEksCluster,
			},
		}),
	}
}

//// LIST FUNCTION

func listEksClusters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := EksClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_cluster.listEksClusters", "get_client_error", err)
		return nil, err
	}

	input := &eks.ListClustersInput{
		MaxResults: aws.Int32(100),
	}

	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < *input.MaxResults {
			if limit < 20 {
				input.MaxResults = aws.Int32(20)
			} else {
				input.MaxResults = aws.Int32(limit)
			}
		}
	}

	paginator := eks.NewListClustersPaginator(svc, input, func(o *eks.ListClustersPaginatorOptions) {
		o.Limit = *input.MaxResults
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_eks_cluster.listEksClusters", "api_error", err)
			return nil, err
		}

		for _, cluster := range output.Clusters {
			d.StreamListItem(ctx, types.Cluster{
				Name: aws.String(cluster),
			})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getEksCluster(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var clusterName string
	if h.Item != nil {
		clusterName = *h.Item.(types.Cluster).Name
	} else {
		clusterName = d.KeyColumnQuals["name"].GetStringValue()
	}

	// create service
	svc, err := EksClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_cluster.getEksCluster", "get_client_error", err)
		return nil, err
	}

	params := &eks.DescribeClusterInput{
		Name: &clusterName,
	}

	op, err := svc.DescribeCluster(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_cluster.getEksCluster", "api_error", err)
		return nil, err
	}

	return op.Cluster, nil
}
