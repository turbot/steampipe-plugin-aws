package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
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
			Hydrate:           getUsagePlan,
		},
		List: &plugin.ListConfig{
			Hydrate: listUsagePlans,
		},
		GetMatrixItem: BuildRegionList,
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

//// LIST FUNCTION

func listUsagePlans(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := APIGatewayService(ctx, d)
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

func getUsagePlan(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getUsagePlan")

	// Create session
	svc, err := APIGatewayService(ctx, d)
	if err != nil {
		return nil, err
	}

	id := d.KeyColumnQuals["id"].GetStringValue()
	params := &apigateway.GetUsagePlanInput{
		UsagePlanId: aws.String(id),
	}

	op, err := svc.GetUsagePlan(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getUsagePlan__", "ERROR", err)
		return nil, err
	}

	return op, nil
}

func getUsagePlanAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getUsagePlanAkas")
	region := d.KeyColumnQualString(matrixKeyRegion)
	usagePlan := h.Item.(*apigateway.UsagePlan)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	akas := []string{"arn:" + commonColumnData.Partition + ":apigateway:" + region + "::/usageplans/" + *usagePlan.Id}

	return akas, nil
}
