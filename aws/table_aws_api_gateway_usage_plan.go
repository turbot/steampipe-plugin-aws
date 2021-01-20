package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAPIGatewayUsagePlan(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_api_gateway_usage_plan",
		Description: "AWS API Gateway Usage Plan",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"NotFoundException"}),
			ItemFromKey:       usagePlanFromKey,
			Hydrate:           getUsagePlan,
		},
		List: &plugin.ListConfig{
			Hydrate: listUsagePlans,
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of a usage plan",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The identifier of a UsagePlan resource",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "product_code",
				Description: "The AWS Markeplace product identifier to associate with the usage plan as a SaaS product on AWS Marketplace",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The description of a usage plan",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "quota",
				Description: "The maximum number of permitted requests per a given unit time interval",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "throttle",
				Description: "The request throttle limits of a usage plan",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "api_stages",
				Description: "The associated API stages of a usage plan",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getUsagePlanAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// ITEM FROM KEY

func usagePlanFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	ID := quals["id"].GetStringValue()
	item := &apigateway.UsagePlan{
		Id: &ID,
	}
	return item, nil
}

//// LIST FUNCTION

func listUsagePlans(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("listUsagePlans", "AWS_REGION", defaultRegion)

	// Create service
	svc, err := APIGatewayService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	err = svc.GetUsagePlansPages(
		&apigateway.GetUsagePlansInput{},
		func(page *apigateway.GetUsagePlansOutput, lastPage bool) bool {
			for _, plan := range page.Items {
				d.StreamListItem(ctx, plan)
			}
			return !lastPage
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getUsagePlan(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getUsagePlan")

	usagePlan := h.Item.(*apigateway.UsagePlan)
	defaultRegion := GetDefaultRegion()

	// Create session
	svc, err := APIGatewayService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &apigateway.GetUsagePlanInput{
		UsagePlanId: usagePlan.Id,
	}

	op, err := svc.GetUsagePlan(params)
	if err != nil {
		logger.Debug("getUsagePlan__", "ERROR", err)
		return nil, err
	}

	return op, nil
}

func getUsagePlanAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getUsagePlanAkas")
	usagePlan := h.Item.(*apigateway.UsagePlan)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	akas := []string{"arn:" + commonColumnData.Partition + ":apigateway:" + commonColumnData.Region + "::/usageplans/" + *usagePlan.Id}

	return akas, nil
}
