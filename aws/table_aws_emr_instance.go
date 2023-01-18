package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/emr"
	"github.com/aws/aws-sdk-go-v2/service/emr/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
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
		GetMatrixItemFunc: BuildRegionList,
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
				Name:        "state",
				Description: "The current state of the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Instance.Status.State"),
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
				Name:        "public_dns_name",
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
				Name:        "state_change_reason",
				Description: "The status change reason details for the instance.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Instance.Status.StateChangeReason"),
			},
			{
				Name:        "status_timeline",
				Description: "The timeline of the instance status over time.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Instance.Status.Timeline"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Instance.Id"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEmrInstanceAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

type emrInstanceInfo struct {
	types.Instance
	ClusterId *string
}

//// LIST FUNCTION

func listEmrInstances(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := EMRClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_emr_instance.listEMRInstances", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Get cluster details
	clusterID := h.Item.(types.ClusterSummary).Id

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

	paginator := emr.NewListInstancesPaginator(svc, input, func(o *emr.ListInstancesPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_emr_cluster.listEmrClusters", "api_error", err)
			return nil, err
		}

		for _, instance := range output.Instances {
			d.StreamListItem(ctx, &emrInstanceInfo{
				Instance:  instance,
				ClusterId: clusterID,
			})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getEmrInstanceAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	data := h.Item.(*emrInstanceInfo)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_emr_instance.getEmrInstanceAkas", "common_data_error", err)
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)

	akas := []string{"arn:" + commonColumnData.Partition + ":emr:" + region + ":" + commonColumnData.AccountId + ":instance/" + *data.Id}

	return akas, nil
}
