package aws

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dax"
	"github.com/aws/aws-sdk-go-v2/service/dax/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsDaxCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dax_cluster",
		Description: "AWS DAX Cluster",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("cluster_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ClusterNotFoundFault", "ServiceLinkedRoleNotFoundFault"}),
			},
			Hydrate: getDaxCluster,
		},
		List: &plugin.ListConfig{
			Hydrate: listDaxClusters,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValueException"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "cluster_name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "cluster_name",
				Description: "The name of the DAX cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) that uniquely identifies the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterArn"),
			},
			{
				Name:        "description",
				Description: "The description of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "active_nodes",
				Description: "The number of nodes in the cluster that are active.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "iam_role_arn",
				Description: "A valid Amazon Resource Name (ARN) that identifies an IAM role.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "node_type",
				Description: "The node type for the nodes in the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "preferred_maintenance_window",
				Description: "A range of time when maintenance of DAX cluster software will be performed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The current status of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subnet_group",
				Description: "The subnet group where the DAX cluster is running.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "total_nodes",
				Description: "The total number of nodes in the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_discovery_endpoint",
				Description: "The configuration endpoint for this DAX cluster, consisting of a DNS name and a port number.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "node_ids_to_remove",
				Description: "A list of nodes to be removed from the cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "nodes",
				Description: "A list of nodes that are currently in the cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "notification_configuration",
				Description: "Describes a notification topic and its status.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "parameter_group",
				Description: "The parameter group being used by nodes in the cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "sse_description",
				Description: "The description of the server-side encryption status on the specified DAX cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SSEDescription"),
			},
			{
				Name:        "security_groups",
				Description: "A list of security groups, and the status of each, for the nodes in the cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags currently associated with the DAX cluster.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDaxClusterTags,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDaxClusterTags,
				Transform:   transform.FromField("Tags").Transform(daxClusterTurbotData),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ClusterArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listDaxClusters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := DAXClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dax_cluster.listDaxClusters", "connection error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
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

	pagesLeft := true
	params := &dax.DescribeClustersInput{
		MaxResults: aws.Int32(maxLimit),
	}

	// Additonal Filter
	equalQuals := d.KeyColumnQuals
	if equalQuals["cluster_name"] != nil {
		params.ClusterNames = []string{equalQuals["cluster_name"].GetStringValue()}
	}

	for pagesLeft {
		result, err := svc.DescribeClusters(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_dax_cluster.listDaxClusters", "api_error", err)
			return nil, err
		}

		for _, cluster := range result.Clusters {
			d.StreamListItem(ctx, cluster)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if result.NextToken != nil {
			pagesLeft = true
			params.NextToken = result.NextToken
		} else {
			pagesLeft = false
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getDaxCluster(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := DAXClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dax_cluster.getDaxCluster", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	quals := d.KeyColumnQuals
	name := quals["cluster_name"].GetStringValue()

	params := &dax.DescribeClustersInput{
		ClusterNames: []string{name},
	}

	op, err := svc.DescribeClusters(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dax_cluster.getDaxCluster", "api_error", err)
		return nil, err
	}

	if op.Clusters != nil && len(op.Clusters) > 0 {
		return op.Clusters[0], nil
	}
	return nil, nil
}

func getDaxClusterTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	clusterArn := *h.Item.(types.Cluster).ClusterArn

	// Create Session
	svc, err := DAXClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dax_cluster.getDaxClusterTags", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &dax.ListTagsInput{
		ResourceName: &clusterArn,
	}

	clusterdata, err := svc.ListTags(ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), "ClusterNotFoundFault") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_dax_cluster.getDaxClusterTags", "api_error", err)
		return nil, err
	}

	return clusterdata, nil
}

//// TRANSFORM FUNCTION

func daxClusterTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*dax.ListTagsOutput)
	if data.Tags == nil {
		return nil, nil
	}

	// Mapping the resource tags inside turbotTags
	if data.Tags != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range data.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
