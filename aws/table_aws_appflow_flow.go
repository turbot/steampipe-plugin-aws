package aws

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/aws/aws-sdk-go/service/cloudcontrolapi"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
)

func tableAwsAppFlowFlow(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_appflow_flow",
		Description: "AWS AppFlow Flow",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("flow_name"),
			Hydrate:    getAppFlowFlow,
		},
		List: &plugin.ListConfig{
			Hydrate: listAppFlowFlows,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "flow_name",
				Description: "Name of the flow.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "flow_arn",
				Description: "ARN identifier of the flow.",
				Hydrate:     getAppFlowFlowARN,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "description",
				Description: "Description of the flow.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAppFlowFlow,
			},
			{
				Name:        "destination_flow_config_list",
				Description: "List of Destination connectors of the flow.",
				Hydrate:     getAppFlowFlow,
				Type:        proto.ColumnType_JSON,
			},
			//{
			//Name:        "kms_arn",
			//Description: "The ARN of the AWS Key Management Service (AWS KMS) key that's used to encrypt your function's environment variables. If it's not provided, AWS Lambda uses a default service key.",
			//Type:        proto.ColumnType_STRING,
			//Hydrate:     getAppFlowFlow,
			//},
			{
				Name:        "source_flow_config",
				Description: "Configurations of Source connector of the flow.",
				Hydrate:     getAppFlowFlow,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tasks",
				Description: "List of tasks for the flow.",
				Hydrate:     getAppFlowFlow,
				Type:        proto.ColumnType_JSON,
			},
			//{
			//Name:        "tags_src",
			//Description: "A list of tags assigned to the flow. NOTE: Missing description.",
			//Type:        proto.ColumnType_JSON,
			//Hydrate:     getAppFlowFlow,
			//Transform:   transform.FromField("Tags"),
			//},
			{
				Name:        "trigger_config",
				Description: "Trigger settings of the flow.",
				Hydrate:     getAppFlowFlow,
				Type:        proto.ColumnType_JSON,
			},

			// Standard Steampipe columns
			//{
			//Name:        "tags",
			//Description: resourceInterfaceDescription("tags"),
			//Type:        proto.ColumnType_JSON,
			//Transform:   transform.FromField("Tags").Transform(getAppFlowFlowTurbotTags),
			//},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAppFlowFlowARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FlowName"),
			},
		}),
	}
}

//// LIST FUNCTION

func listAppFlowFlows(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAppFlowFlows")

	// Create session
	svc, err := CloudControlAPIService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := cloudcontrolapi.ListResourcesInput{TypeName: types.String("AWS::AppFlow::Flow")}

	err = svc.ListResourcesPages(&input,
		func(page *cloudcontrolapi.ListResourcesOutput, isLast bool) bool {
			for _, flow := range page.ResourceDescriptions {
				properties := flow.Properties
				var jsonMap map[string]interface{}
				json.Unmarshal([]byte(*properties), &jsonMap)

				d.StreamListItem(ctx, jsonMap)
				// This will return zero if context has been cancelled (i.e due to manual cancellation) or
				// if there is a limit, it will return the number of rows required to reach this limit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return true
				}
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAppFlowFlow(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAppFlowFlow")

	// Create session
	svc, err := CloudControlAPIService(ctx, d)
	if err != nil {
		return nil, err
	}

	var identifier string

	if h.Item != nil {
		result := h.Item
		reflectedValue := reflect.ValueOf(result).MapIndex(reflect.ValueOf("FlowName"))
		identifier = reflectedValue.Interface().(string)
	} else {
		identifier = d.KeyColumnQuals["flow_name"].GetStringValue()
	}

	input := &cloudcontrolapi.GetResourceInput{
		Identifier: types.String(identifier),
		TypeName:   types.String("AWS::AppFlow::Flow"),
	}

	item, err := svc.GetResource(input)
	if err != nil {
		return nil, err
	}

	properties := item.ResourceDescription.Properties
	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(*properties), &jsonMap)

	return jsonMap, nil
}

//// TRANSFORM FUNCTIONS

func getAppFlowFlowARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAppFlowFlowARN")

	if h.Item == nil {
		return nil, nil
	}

	var flowARN string

	result := h.Item
	reflectedARN := reflect.ValueOf(result).MapIndex(reflect.ValueOf("FlowArn"))

	// FlowArn property exists when listing flows
	if reflectedARN.IsValid() {
		flowARN = reflectedARN.Interface().(string)
		return flowARN, nil
	}

	flowName := d.KeyColumnQuals["flow_name"].GetStringValue()
	region := d.KeyColumnQualString(matrixKeyRegion)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("getEc2InstanceARN", "getCommonColumnsCached_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	flowARN = "arn:" + commonColumnData.Partition + ":appflow:" + region + ":" + commonColumnData.AccountId + ":flow/" + flowName

	return flowARN, nil
}
