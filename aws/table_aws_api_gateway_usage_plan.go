package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"

	apigatewayv1 "github.com/aws/aws-sdk-go/service/apigateway"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAPIGatewayUsagePlan(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_api_gateway_usage_plan",
		Description: "AWS API Gateway Usage Plan",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException"}),
			},
			Hydrate: getUsagePlan,
			Tags:    map[string]string{"service": "apigateway", "action": "GetUsagePlan"},
		},
		List: &plugin.ListConfig{
			Hydrate: listUsagePlans,
			Tags:    map[string]string{"service": "apigateway", "action": "GetUsagePlans"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(apigatewayv1.EndpointsID),
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
	logger := plugin.Logger(ctx)

	// Create service
	svc, err := APIGatewayClient(ctx, d)
	if err != nil {
		logger.Error("aws_api_gateway_usage_plan.listUsagePlans", "service_client_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(500)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	input := &apigateway.GetUsagePlansInput{
		Limit: aws.Int32(maxLimit),
	}

	paginator := apigateway.NewGetUsagePlansPaginator(svc, input, func(o *apigateway.GetUsagePlansPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_api_gateway_rest_api.listUsagePlans", "api_error", err)
			return nil, err
		}

		for _, items := range output.Items {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getUsagePlan(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := APIGatewayClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gateway_rest_api.getUsagePlan", "service_cllient_error", err)
		return nil, err
	}

	id := d.EqualsQuals["id"].GetStringValue()
	params := &apigateway.GetUsagePlanInput{
		UsagePlanId: aws.String(id),
	}

	op, err := svc.GetUsagePlan(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gateway_usage_plan.getUsagePlan", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getUsagePlanAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	id := ""

	switch h.Item.(type) {
	case *apigateway.GetUsagePlanOutput:
		id = *h.Item.(*apigateway.GetUsagePlanOutput).Id
	case types.UsagePlan:
		id = *h.Item.(types.UsagePlan).Id
	}

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	akas := []string{"arn:" + commonColumnData.Partition + ":apigateway:" + region + "::/usageplans/" + id}

	return akas, nil
}
