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

func tableAwsEksNodeGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_eks_node_group",
		Description: "AWS EKS Node Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"nodegroup_name", "cluster_name"}),
			Hydrate:    getEKSNodeGroup,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundErrorV2([]string{"InvalidParameterException", "InvalidParameter"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listEKSClusters,
			Hydrate:       listEKSNodeGroups,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "cluster_name", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "nodegroup_name",
				Description: "The name associated with an Amazon EKS managed node group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) associated with the managed node group.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEKSNodeGroup,
				Transform:   transform.FromField("NodegroupArn"),
			},
			{
				Name:        "cluster_name",
				Description: "The name of the cluster that the managed node group resides in.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at",
				Description: "The Unix epoch timestamp in seconds for when the managed node group was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getEKSNodeGroup,
			},
			{
				Name:        "status",
				Description: "The current status of the managed node group.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEKSNodeGroup,
			},
			{
				Name:        "ami_type",
				Description: "The AMI type that was specified in the node group configuration.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEKSNodeGroup,
			},
			{
				Name:        "capacity_type",
				Description: "The capacity type of your managed node group.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEKSNodeGroup,
			},
			{
				Name:        "disk_size",
				Description: "The disk size in the node group configuration.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getEKSNodeGroup,
			},
			{
				Name:        "modified_at",
				Description: "The Unix epoch timestamp in seconds for when the managed node group was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getEKSNodeGroup,
			},
			{
				Name:        "node_role",
				Description: "The IAM role associated with your node group.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEKSNodeGroup,
			},
			{
				Name:        "release_version",
				Description: "If the node group was deployed using a launch template with a custom AMI, then this is the AMI ID that was specified in the launch template. For node groups that weren't deployed using a launch template, this is the version of the Amazon EKS optimized AMI that the node group was deployed with.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEKSNodeGroup,
			},
			{
				Name:        "version",
				Description: "The Kubernetes version of the managed node group.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEKSNodeGroup,
			},
			{
				Name:        "health",
				Description: "The health status of the node group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEKSNodeGroup,
			},
			{
				Name:        "instance_types",
				Description: "The instance type that is associated with the node group. If the node group was deployed with a launch template, then this is null.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEKSNodeGroup,
			},
			{
				Name:        "labels",
				Description: "The Kubernetes labels applied to the nodes in the node group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEKSNodeGroup,
			},
			{
				Name:        "launch_template",
				Description: "If a launch template was used to create the node group, then this is the launch template that was used.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEKSNodeGroup,
			},
			{
				Name:        "remote_access",
				Description: "The remote access configuration that is associated with the node group. If the node group was deployed with a launch template, then this is null.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEKSNodeGroup,
			},
			{
				Name:        "resources",
				Description: "The resources associated with the node group, such as Auto Scaling groups and security groups for remote access.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEKSNodeGroup,
			},
			{
				Name:        "scaling_config",
				Description: "The scaling configuration details for the Auto Scaling group that is associated with your node group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEKSNodeGroup,
			},
			{
				Name:        "subnets",
				Description: "The subnets that were specified for the Auto Scaling group that is associated with your node group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEKSNodeGroup,
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEKSNodeGroup,
			},
			{
				Name:        "taints",
				Description: "The Kubernetes taints to be applied to the nodes in the node group when they are created.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEKSNodeGroup,
			},
			{
				Name:        "update_config",
				Description: "The node group update configuration.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEKSNodeGroup,
			},

			// Steampipe Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("NodegroupName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("NodegroupArn").Transform(transform.EnsureStringArray),
				Hydrate:     getEKSNodeGroup,
			},
		}),
	}
}

//// LIST FUNCTION

func listEKSNodeGroups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get cluster details
	cluster := d.KeyColumnQuals["cluster_name"].GetStringValue()
	clusterName := *h.Item.(types.Cluster).Name

	if cluster != "" {
		if cluster != clusterName {
			return nil, nil
		}
	}

	// Create service
	svc, err := EKSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_node_group.listEksNodeGroups", "connection_error", err)
		return nil, err
	}

	maxItems := int32(100)
	input := &eks.ListNodegroupsInput{
		ClusterName: &clusterName,
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	input.MaxResults = aws.Int32(maxItems)
	paginator := eks.NewListNodegroupsPaginator(svc, input, func(o *eks.ListNodegroupsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_eks_node_group.listEKSNodeGroups", "api_error", err)
			return nil, err
		}

		for _, nodegroup := range output.Nodegroups {
			d.StreamListItem(ctx, &types.Nodegroup{
				NodegroupName: aws.String(nodegroup),
				ClusterName:   aws.String(clusterName),
			})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getEKSNodeGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var clusterName, nodegroupName string
	if h.Item != nil {
		clusterName = *h.Item.(*types.Nodegroup).ClusterName
		nodegroupName = *h.Item.(*types.Nodegroup).NodegroupName
	} else {
		clusterName = d.KeyColumnQuals["cluster_name"].GetStringValue()
		nodegroupName = d.KeyColumnQuals["nodegroup_name"].GetStringValue()
	}

	// check if clusterName or nodegroupName is empty
	if clusterName == "" || nodegroupName == "" {
		return nil, nil
	}

	// create service
	svc, err := EKSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_node_group.getEKSNodeGroup", "connection_error", err)
		return nil, err
	}

	params := &eks.DescribeNodegroupInput{
		ClusterName:   &clusterName,
		NodegroupName: &nodegroupName,
	}

	op, err := svc.DescribeNodegroup(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_eks_node_group.getEKSNodeGroup", "api_error", err)
		return nil, err
	}

	return op.Nodegroup, nil
}
