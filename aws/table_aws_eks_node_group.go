package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/service/eks"
)

//// TABLE DEFINITION

func tableAwsEksNodeGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_eks_node_group",
		Description: "AWS Elastic Kubernetes Service Node Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "cluster_name"}),
			Hydrate:    getEksNodeGroup,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listEksClusters,
			Hydrate:       listEksNodeGroup,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name associated with an Amazon EKS managed node group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("NodegroupName"),
			},
			{
				Name:        "node_group_arn",
				Description: "The Amazon Resource Name (ARN) of the cluster.",
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
				Name:        "version",
				Description: "The Kubernetes version of the managed node group.",
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
				Name:        "created_at",
				Description: "The Unix epoch timestamp in seconds for when the managed node group was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "modified_at",
				Description: "The timestamp when the managed node group was last modified.",
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
				Name:        "capacity_type",
				Description: "The capacity type of the managed node group.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "ami_type",
				Description: "If the node group was deployed using a launch template with a custom AMI, then this is CUSTOM. For node groups that weren't deployed using a launch template, this is the AMI type that was specified in the node group configuration.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "node_role",
				Description: "The IAM role associated with your node group.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "disk_size",
				Description: "If the node group wasn't deployed with a launch template, then this is the disk size in the node group configuration. If the node group was deployed with a launch template, then this is null.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "scaling_config",
				Description: "The scaling configuration details for the Auto Scaling group that is associated with node group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "instance_types",
				Description: "If the node group wasn't deployed with a launch template, then this is the instance type that is associated with the node group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "subnets",
				Description: "The subnets that were specified for the Auto Scaling group that is associated with node group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "remote_access",
				Description: "If the node group wasn't deployed with a launch template, then this is the remote access configuration that is associated with the node group.",
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
				Name:        "resources",
				Description: "The resources associated with the node group, such as Auto Scaling groups and security groups for remote access.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "health",
				Description: "The health status of the node group.",
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
				Name:        "tags",
				Description: "The metadata applied to the node group to assist with categorization and organization.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEksNodeGroup,
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("NodegroupName"),
				Hydrate:     getEksNodeGroup,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("NodegroupArn").Transform(arnToAkas),
				Hydrate:     getEksNodeGroup,
			},
		}),
	}
}

//// LIST FUNCTION

func listEksNodeGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listEksClusters", "AWS_REGION", region)

	// Create service
	svc, err := EksService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	clusterName := *h.Item.(*eks.Cluster).Name

	err = svc.ListNodegroupsPages(
		&eks.ListNodegroupsInput{
			ClusterName: &clusterName,
		},
		func(page *eks.ListNodegroupsOutput, isLast bool) bool {
			for _, nodeGroup := range page.Nodegroups {
				d.StreamLeafListItem(ctx, &eks.Nodegroup{
					NodegroupName: nodeGroup,
					ClusterName:   &clusterName,
				})
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getEksNodeGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEksCluster")
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	var clusterName string
	var nodeGroupName string
	if h.Item != nil {
		clusterName = *h.Item.(*eks.Nodegroup).ClusterName
		nodeGroupName = *h.Item.(*eks.Nodegroup).NodegroupName
	} else {
		clusterName = d.KeyColumnQuals["cluster_name"].GetStringValue()
		nodeGroupName = d.KeyColumnQuals["name"].GetStringValue()
	}

	// create service
	svc, err := EksService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &eks.DescribeNodegroupInput{
		ClusterName:   &clusterName,
		NodegroupName: &nodeGroupName,
	}

	op, err := svc.DescribeNodegroup(params)
	if err != nil {
		return nil, err
	}

	return op.Nodegroup, nil
}
