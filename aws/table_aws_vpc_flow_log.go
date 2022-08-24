package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsVpcFlowlog(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_flow_log",
		Description: "AWS VPC Flowlog",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("flow_log_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"Client.InvalidInstanceID.NotFound", "InvalidParameterValue"}),
			},
			Hydrate: getVpcFlowlog,
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcFlowlogs,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "deliver_logs_status", Require: plugin.Optional},
				{Name: "log_destination_type", Require: plugin.Optional},
				{Name: "log_group_name", Require: plugin.Optional},
				{Name: "resource_id", Require: plugin.Optional},
				{Name: "traffic_type", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "flow_log_id",
				Description: "The ID of the flow log.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The date and time the flow log was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "deliver_logs_error_message",
				Description: "Information about the error that occurred.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "deliver_logs_permission_arn",
				Description: "The ARN of the IAM role that posts logs to CloudWatch Logs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "deliver_logs_status",
				Description: "The status of the logs delivery (SUCCESS | FAILED).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "flow_log_status",
				Description: "The status of the flow log (ACTIVE).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "log_group_name",
				Description: "The name of the flow log group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_id",
				Description: "The ID of the VPC, subnet, or network interface.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "traffic_type",
				Description: "The type of traffic. Valid values are: 'ACCEPT', 'REJECT',  'ALL'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "log_destination_type",
				Description: "Specifies the type of destination to which the flow log data is published.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "log_destination",
				Description: "Specifies the destination to which the flow log data is published.",
				Type:        proto.ColumnType_STRING,
			},
			// bucket_name is a simpler view of the log destination bucket name, without the full path
			{
				Name:        "bucket_name",
				Description: "The name of the destination bucket to which the flow log data is published.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(logDestinationBucketName),
			},
			{
				Name:        "log_format",
				Description: "The format of the flow log record.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "max_aggregation_interval",
				Description: "The maximum interval of time, in seconds, during which a flow of packets is captured and aggregated into a flow log record.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the VPC flowlog.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			//standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(vpcFlowlogTurbotData, "Tags"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(vpcFlowlogTurbotData, "Title"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVpcFlowlogAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listVpcFlowlogs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	plugin.Logger(ctx).Trace("listVpcFlowlogs", "AWS_REGION", region)

	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// The max page limit is not mentioned in the doc, so here the max limt set to 1000 and min to 1
	// https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeFlowLogs.html
	input := &ec2.DescribeFlowLogsInput{
		MaxResults: aws.Int64(1000),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
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

	filterKeyMap := []VpcFilterKeyMap{
		{ColumnName: "deliver_logs_status", FilterName: "deliver-log-status", ColumnType: "string"},
		{ColumnName: "log_destination_type", FilterName: "log-destination-type", ColumnType: "string"},
		{ColumnName: "log_group_name", FilterName: "log-group-name", ColumnType: "string"},
		{ColumnName: "resource_id", FilterName: "resource-id", ColumnType: "string"},
		{ColumnName: "traffic_type", FilterName: "traffic-type", ColumnType: "string"},
	}

	filters := buildVpcResourcesFilterParameter(filterKeyMap, d.Quals)
	if len(filters) > 0 {
		input.Filter = filters
	}

	err = svc.DescribeFlowLogsPages(
		input,
		func(page *ec2.DescribeFlowLogsOutput, lastPage bool) bool {
			for _, item := range page.FlowLogs {
				d.StreamListItem(ctx, item)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !lastPage
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getVpcFlowlog(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getVpcFlowlog")

	quals := d.KeyColumnQuals
	flowlogID := quals["flow_log_id"].GetStringValue()

	region := d.KeyColumnQualString(matrixKeyRegion)

	// Create session
	svc, err := Ec2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	params := &ec2.DescribeFlowLogsInput{
		FlowLogIds: []*string{&flowlogID},
	}

	//get call
	item, err := svc.DescribeFlowLogs(params)
	if err != nil {
		logger.Debug("getVpcFlowlogs__", "Error", err)
		return nil, err
	}

	if item.FlowLogs != nil && len(item.FlowLogs) > 0 {
		return item.FlowLogs[0], nil
	}

	return nil, nil
}

func getVpcFlowlogAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcFlowlogAkas")
	region := d.KeyColumnQualString(matrixKeyRegion)
	vpcFlowlog := h.Item.(*ec2.FlowLog)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":vpc-flow-log/" + *vpcFlowlog.FlowLogId}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func vpcFlowlogTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	vpcFlowlog := d.HydrateItem.(*ec2.FlowLog)
	param := d.Param.(string)

	// Get resource title
	title := vpcFlowlog.FlowLogId

	// Get the resource tags
	var turbotTagsMap map[string]string
	if vpcFlowlog.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range vpcFlowlog.Tags {
			turbotTagsMap[*i.Key] = *i.Value
			if *i.Key == "Name" {
				title = i.Value
			}
		}
	}

	if param == "Tags" {
		return turbotTagsMap, nil
	}

	return title, nil
}

func logDestinationBucketName(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*ec2.FlowLog)
	logDestination := types.SafeString(data.LogDestination)
	if logDestination == "" {
		return nil, nil
	}
	splitData := strings.Split(logDestination, ":")
	return splitData[len(splitData)-1], nil
}
