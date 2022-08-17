package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/emr"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

//// TABLE DEFINITION

func tableAwsEmrInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_emr_instance",
		Description: "AWS EMR Instance",
		List: &plugin.ListConfig{
			ParentHydrate: listEmrClusters,
			Hydrate:       listEmrInstances,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "cluster_id", Require: plugin.Optional},
				{Name: "instance_fleet_id", Require: plugin.Optional},
				{Name: "instance_group_id", Require: plugin.Optional},
			},
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The unique identifier for the instance in Amazon EMR.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Instance.Id"),
			},
			{
				Name:        "cluster_id",
				Description: "The unique identifier for the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ec2_instance_id",
				Description: "The unique identifier of the instance in Amazon EC2.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Instance.Ec2InstanceId"),
			},
			{
				Name:        "instance_fleet_id",
				Description: "The unique identifier of the instance fleet to which an EC2 instance belongs.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Instance.InstanceFleetId"),
			},
			{
				Name:        "instance_group_id",
				Description: "The identifier of the instance group to which this instance belongs.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Instance.InstanceGroupId"),
			},
			{
				Name:        "instance_type",
				Description: "The EC2 instance type, for example m3.xlarge.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Instance.InstanceType"),
			},
			{
				Name:        "market",
				Description: "The instance purchasing option. Valid values are ON_DEMAND or SPOT.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Instance.Market"),
			},
			{
				Name:        "private_dns_name",
				Description: "The private DNS name of the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Instance.PrivateDnsName"),
			},
			{
				Name:        "private_ip_address",
				Description: "The private IP address of the instance.",
				Type:        proto.ColumnType_IPADDR,
				Transform:   transform.FromField("Instance.PrivateIpAddress"),
			},
			{
				Name:        "PublicDnsName",
				Description: "The public DNS name of the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Instance.PublicDnsName"),
			},
			{
				Name:        "public_ip_address",
				Description: "The public IP address of the instance.",
				Type:        proto.ColumnType_IPADDR,
				Transform:   transform.FromField("Instance.PublicIpAddress"),
			},
			{
				Name:        "ebs_volumes",
				Description: "The list of Amazon EBS volumes that are attached to this instance.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Instance.EbsVolumes"),
			},
			{
				Name:        "status",
				Description: "The current status of the instance.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Instance.Status"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},
		}),
	}
}

type EmrInstanceInfo struct {
	*emr.Instance
	ClusterId *string
}

//// LIST FUNCTION

func listEmrInstances(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := EmrService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get cluster details
	clusterID := h.Item.(*emr.ClusterSummary).Id

	if d.KeyColumnQualString("cluster_id") != "" && d.KeyColumnQualString("cluster_id") != *clusterID {
		return nil, nil
	}

	if d.KeyColumnQualString("cluster_id") != "" {
		clusterID = aws.String(d.KeyColumnQualString("cluster_id"))
	}

	input := &emr.ListInstancesInput{
		ClusterId: clusterID,
	}

	if d.KeyColumnQualString("instance_fleet_id") != "" {
		input.InstanceFleetId = aws.String(d.KeyColumnQualString("instance_fleet_id"))
	}

	if d.KeyColumnQualString("instance_group_id") != "" {
		input.InstanceGroupId = aws.String(d.KeyColumnQualString("instance_group_id"))
	}

	// List call
	err = svc.ListInstancesPages(
		input,
		func(page *emr.ListInstancesOutput, isLast bool) bool {
			for _, instance := range page.Instances {
				d.StreamListItem(ctx, &EmrInstanceInfo{
					Instance:  instance,
					ClusterId: clusterID,
				})

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	return nil, err
}
