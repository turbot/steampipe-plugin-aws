package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"

	ssmv1 "github.com/aws/aws-sdk-go/service/ssm"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSSMManagedInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssm_managed_instance",
		Description: "AWS SSM Managed Instance",
		List: &plugin.ListConfig{
			Hydrate: listSsmManagedInstances,
			Tags:    map[string]string{"service": "ssm", "action": "DescribeInstanceInformation"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "instance_id", Require: plugin.Optional},
				{Name: "agent_version", Require: plugin.Optional},
				{Name: "ping_status", Require: plugin.Optional},
				{Name: "platform_type", Require: plugin.Optional},
				{Name: "activation_id", Require: plugin.Optional},
				{Name: "resource_type", Require: plugin.Optional},
				{Name: "association_status", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ssmv1.EndpointsID),
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
				Type:        proto.ColumnType_STRING,
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

	// Create session
	svc, err := SSMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_managed_instance.listSsmManagedInstances", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	maxItems := int32(50)
	input := &ssm.DescribeInstanceInformationInput{}

	filters := buildSSMManagedInstanceFilter(d.Quals)
	if len(filters) > 0 {
		input.Filters = filters
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 5 {
				maxItems = int32(5)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	input.MaxResults = aws.Int32(maxItems)
	paginator := ssm.NewDescribeInstanceInformationPaginator(svc, input, func(o *ssm.DescribeInstanceInformationPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ssm_managed_instance.listSsmManagedInstances", "api_error", err)
			return nil, err
		}

		for _, managedInstance := range output.InstanceInformationList {
			d.StreamListItem(ctx, managedInstance)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSsmManagedInstanceARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(types.InstanceInformation)
	region := d.EqualsQualString(matrixKeyRegion)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_managed_instance.getSsmManagedInstanceARN", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	arn := "arn:" + commonColumnData.Partition + ":ssm:" + region + ":" + commonColumnData.AccountId + ":managed-instance/" + *data.InstanceId

	return arn, nil
}

//// UTILITY FUNCTION

// Build ssm managed instance list call input filter
func buildSSMManagedInstanceFilter(quals plugin.KeyColumnQualMap) []types.InstanceInformationStringFilter {
	filters := make([]types.InstanceInformationStringFilter, 0)

	filterQuals := map[string]string{
		"instance_id":        "InstanceIds",
		"agent_version":      "AgentVersion",
		"ping_status":        "PingStatus",
		"platform_type":      "PlatformTypes",
		"activation_id":      "ActivationIds",
		"resource_type":      "ResourceType",
		"association_status": "AssociationStatus",
	}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := types.InstanceInformationStringFilter{
				Key: aws.String(filterName),
			}

			value := getQualsValueByColumn(quals, columnName, "string")
			val, ok := value.(string)
			if ok {
				filter.Values = []string{val}
			} else {
				filter.Values = value.([]string)
			}
			filters = append(filters, filter)
		}
	}
	return filters
}
