package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
)

//// TABLE DEFINITION

func tableAwsEksNodeGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_eks_node_group",
		Description: "AWS EKS Node Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"nodegroup_name", "cluster_name"}),
			Hydrate:    getEksNodeGroup,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidParameterException", "InvalidParameter"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listEksClusters,
			Hydrate:       listEksNodeGroups,
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
				Hydrate:     getEksNodeGroup,
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
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "status",
				Description: "The current status of the managed node group.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "ami_type",
				Description: "The AMI type that was specified in the node group configuration.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "capacity_type",
				Description: "The capacity type of your managed node group.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "disk_size",
				Description: "The disk size in the node group configuration.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "modified_at",
				Description: "The Unix epoch timestamp in seconds for when the managed node group was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "node_role",
				Description: "The IAM role associated with your node group.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "release_version",
				Description: "If the node group was deployed using a launch template with a custom AMI, then this is the AMI ID that was specified in the launch template. For node groups that weren't deployed using a launch template, this is the version of the Amazon EKS optimized AMI that the node group was deployed with.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "version",
				Description: "The Kubernetes version of the managed node group.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "health",
				Description: "The health status of the node group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "instance_types",
				Description: "The instance type that is associated with the node group. If the node group was deployed with a launch template, then this is null.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "labels",
				Description: "The Kubernetes labels applied to the nodes in the node group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "launch_template",
				Description: "If a launch template was used to create the node group, then this is the launch template that was used.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "remote_access",
				Description: "The remote access configuration that is associated with the node group. If the node group was deployed with a launch template, then this is null.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "resources",
				Description: "The resources associated with the node group, such as Auto Scaling groups and security groups for remote access.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "scaling_config",
				Description: "The scaling configuration details for the Auto Scaling group that is associated with your node group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "subnets",
				Description: "The subnets that were specified for the Auto Scaling group that is associated with your node group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "taints",
				Description: "The Kubernetes taints to be applied to the nodes in the node group when they are created.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "update_config",
				Description: "The node group update configuration.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksNodeGroup,
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
				Hydrate:     getEksNodeGroup,
			},
		}),
	}
}

//// LIST FUNCTION

func listEksNodeGroups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get cluster details
	cluster := d.KeyColumnQuals["cluster_name"].GetStringValue()
	clusterName := *h.Item.(*eks.Cluster).Name

	if cluster != "" {
		if cluster != clusterName {
			return nil, nil
		}
	}

	// Create service
	svc, err := EksService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &eks.ListNodegroupsInput{
		ClusterName: &clusterName,
		MaxResults:  aws.Int64(100),
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 1 {
				input.MaxResults = aws.Int64(1)
			} else {
				input.MaxResults = limit
			}
		}
	}

	err = svc.ListNodegroupsPages(
		input,
		func(page *eks.ListNodegroupsOutput, _ bool) bool {
			for _, nodegroup := range page.Nodegroups {
				d.StreamListItem(ctx, &eks.Nodegroup{
					NodegroupName: nodegroup,
					ClusterName:   &clusterName,
				})
				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return true
		},
	)

	if err != nil {
		plugin.Logger(ctx).Error("ListNodegroupsPages", "api_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getEksNodeGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Error("getEksNodeGroup")

	var clusterName, nodegroupName string
	if h.Item != nil {
		clusterName = *h.Item.(*eks.Nodegroup).ClusterName
		nodegroupName = *h.Item.(*eks.Nodegroup).NodegroupName
	} else {
		clusterName = d.KeyColumnQuals["cluster_name"].GetStringValue()
		nodegroupName = d.KeyColumnQuals["nodegroup_name"].GetStringValue()
	}

	// check if clusterName or nodegroupName is empty
	if clusterName == "" || nodegroupName == "" {
		return nil, nil
	}

	// create service
	svc, err := EksService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &eks.DescribeNodegroupInput{
		ClusterName:   &clusterName,
		NodegroupName: &nodegroupName,
	}

	op, err := svc.DescribeNodegroup(params)

	if err != nil {
		plugin.Logger(ctx).Error("DescribeNodegroup", "api_error", err)
		return nil, err
	}

	return op.Nodegroup, nil
}
