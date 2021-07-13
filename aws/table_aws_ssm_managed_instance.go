package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSSMManagedInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssm_managed_instance",
		Description: "AWS SSM Managed Instance",
		List: &plugin.ListConfig{
			Hydrate: listSsmManagedInstances,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name assigned to an on-premises server or virtual machine (VM) when it is activated as a Systems Manager managed instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_id",
				Description: "The ID of the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the instance.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSsmManagedInstanceARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "resource_type",
				Description: "The type of instance. Instances are either EC2 instances or managed instances.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ping_status",
				Description: "Connection status of SSM Agent.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "activation_id",
				Description: "The activation ID created by Systems Manager when the server or VM was registered.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "agent_version",
				Description: "The version of SSM Agent running on your Linux instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "association_status",
				Description: "The status of the association.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "computer_name",
				Description: "The fully qualified host name of the managed instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ip_address",
				Description: "The IP address of the managed instance.",
				Type:        proto.ColumnType_IPADDR,
				Transform:   transform.FromField("IPAddress"),
			},
			{
				Name:        "is_latest_version",
				Description: "Indicates whether the latest version of SSM Agent is running on your Linux Managed Instance.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "last_association_execution_date",
				Description: "The date the association was last run.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_ping_date_time",
				Description: "The date and time when the agent last pinged the Systems Manager service.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_successful_association_execution_date",
				Description: "The last date the association was successfully run.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "platform_name",
				Description: "The name of the operating system platform running on your instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "platform_type",
				Description: "The operating system platform type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "platform_version",
				Description: "The version of the OS platform running on your instance.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "registration_date",
				Description: "The date the server or VM was registered with AWS as a managed instance.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "association_overview",
				Description: "Information about the association.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InstanceId"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSsmManagedInstanceARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listSsmManagedInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listSsmManagedInstances")

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listSsmManagedInstances", "AWS_REGION", region)

	// Create session
	svc, err := SsmService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.DescribeInstanceInformationPages(
		&ssm.DescribeInstanceInformationInput{},
		func(page *ssm.DescribeInstanceInformationOutput, isLast bool) bool {
			for _, managedInstance := range page.InstanceInformationList {
				d.StreamListItem(ctx, managedInstance)
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getSsmManagedInstanceARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSsmManagedInstanceARN")
	data := h.Item.(*ssm.InstanceInformation)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	arn := "arn:" + commonColumnData.Partition + ":ssm:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":managed-instance/" + *data.InstanceId

	return arn, nil
}
