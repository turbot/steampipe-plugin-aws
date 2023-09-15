package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"

	apigatewayv1 "github.com/aws/aws-sdk-go/service/apigatewayv2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsAPIGatewayDomainName(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_api_gateway_domain_name",
		Description: "AWS API Gateway Domain Name",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"domain_name"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException"}),
			},
			Hydrate: getApiGatewayDomainName,
			Tags:    map[string]string{"service": "apigateway", "action": "GetDomainName"},
		},
		List: &plugin.ListConfig{
			Hydrate: listApiGatewayDomainNames,
			Tags:    map[string]string{"service": "apigateway", "action": "GetDomainNames"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(apigatewayv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "domain_name",
				Description: "The custom domain name as an API host name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "certificate_arn",
				Description: "The reference to an AWS-managed certificate that will be used by edge-optimized endpoint for this domain name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "certificate_name",
				Description: "The name of the certificate that will be used by edge-optimized endpoint for this domain name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "certificate_upload_date",
				Description: "The timestamp when the certificate that was used by edge-optimized endpoint for this domain name was uploaded.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "distribution_domain_name",
				Description: "The domain name of the Amazon CloudFront distribution associated with this custom domain name for an edge-optimized endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "distribution_hosted_zone_id",
				Description: "The region-agnostic Amazon Route 53 Hosted Zone ID of the edge-optimized endpoint. The valid value is Z2FDTNDATAQYW2 for all the regions.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "domain_name_status",
				Description: "The status of the DomainName migration. The valid values are AVAILABLE and UPDATING. If the status is UPDATING, the domain cannot be modified further until the existing operation is complete.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "domain_name_status_message",
				Description: "An optional text message containing detailed information about status of the DomainName migration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ownership_verification_certificate_arn",
				Description: "The ARN of the public certificate issued by ACM to validate ownership of your custom domain. Only required when configuring mutual TLS and using an ACM imported or private CA certificate ARN as the regionalCertificateArn.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "regional_certificate_arn",
				Description: "The reference to an AWS-managed certificate that will be used for validating the regional domain name. AWS Certificate Manager is the only supported source.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "regional_certificate_name",
				Description: "The name of the certificate that will be used for validating the regional domain name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "regional_domain_name",
				Description: "The domain name associated with the regional endpoint for this custom domain name. You set up this association by adding a DNS record that points the custom domain name to this regional domain name. The regional domain name is returned by API Gateway when you create a regional endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "regional_hosted_zone_id",
				Description: "The region-specific Amazon Route 53 Hosted Zone ID of the regional endpoint. For more information, see Set up a Regional Custom Domain Name and AWS Regions and Endpoints for API Gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "security_policy",
				Description: "The Transport Layer Security (TLS) version + cipher suite for this DomainName. The valid values are TLS_1_0 and TLS_1_2.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "endpoint_configuration",
				Description: "The endpoint configuration of this DomainName showing the endpoint types of the domain name.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "mutual_tls_authentication",
				Description: "The mutual TLS authentication configuration for a custom domain name. If specified, API Gateway performs two-way authentication between the client and the server. Clients must present a trusted certificate to access your API.",
				Type:        proto.ColumnType_JSON,
			},

			//// Steampipe standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DomainName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getapiGatewayDomainNameAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listApiGatewayDomainNames(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := APIGatewayClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gateway_domain_name.listApiGatewayDomainNames", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(500)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	params := &apigateway.GetDomainNamesInput{
		Limit: aws.Int32(maxLimit),
	}

	paginator := apigateway.NewGetDomainNamesPaginator(svc, params, func(o *apigateway.GetDomainNamesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)
		
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_api_gateway_domain_name.listApiGatewayDomainNames", "api_error", err)
			return nil, err
		}

		for _, domain := range output.Items {
			d.StreamListItem(ctx, domain)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getApiGatewayDomainName(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := APIGatewayClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gateway_domain_name.getApiGatewayDomainName", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	domainName := d.EqualsQuals["domain_name"].GetStringValue()

	// Empty check
	if domainName == "" {
		return nil, nil
	}

	input := &apigateway.GetDomainNameInput{
		DomainName: aws.String(domainName),
	}

	op, err := svc.GetDomainName(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gateway_domain_name.getApiGatewayDomainName", "api_error", err)
		return nil, err
	}

	if op != nil {
		return op, nil
	}

	return nil, nil
}

func getapiGatewayDomainNameAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	domainName := ""

	switch h.Item.(type) {
	case *types.DomainName:
		domainName = *h.Item.(*types.DomainName).DomainName
	case *apigateway.GetDomainNameOutput:
		domainName = *h.Item.(*apigateway.GetDomainNameOutput).DomainName
	}

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_api_gateway_domain_name.getapiGatewayDomainNameAkas", "getCommonColumns_error", err)
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)
	akas := []string{"arn:" + commonColumnData.Partition + ":apigateway:" + region + "::/domainname/" + domainName}

	return akas, nil
}
