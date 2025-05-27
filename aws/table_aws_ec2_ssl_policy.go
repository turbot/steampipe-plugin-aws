package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEc2SslPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_ssl_policy",
		Description: "AWS EC2 SSL Policy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "region"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"SSLPolicyNotFound"}),
			},
			Hydrate: getEc2SslPolicy,
			Tags:    map[string]string{"service": "elasticloadbalancing", "action": "DescribeSSLPolicies"},
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2SslPolicies,
			Tags:    map[string]string{"service": "elasticloadbalancing", "action": "DescribeSSLPolicies"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_ELASTICLOADBALANCING_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ciphers",
				Description: "A list of ciphers.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ssl_protocols",
				Description: "A list of protocols.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "supported_load_balancer_types",
				Description: "The supported load balancers.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
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
				Hydrate:     getEc2SslPolicyAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listEc2SslPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := ELBV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_ssl_policy.listEc2SslPolicies", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(400)
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

	// List call
	params := &elasticloadbalancingv2.DescribeSSLPoliciesInput{
		PageSize: aws.Int32(maxLimit),
	}

	pagesLeft := true
	for pagesLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		response, err := svc.DescribeSSLPolicies(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ec2_ssl_policy.listEc2SslPolicies", "api_error", err)
			return nil, err
		}

		for _, sslPolicy := range response.SslPolicies {
			d.StreamListItem(ctx, sslPolicy)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if response.NextMarker != nil {
			pagesLeft = true
			params.Marker = response.NextMarker
		} else {
			pagesLeft = false
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getEc2SslPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	matrixKeyRegion := d.EqualsQualString(matrixKeyRegion)
	name := d.EqualsQuals["name"].GetStringValue()
	regionName := d.EqualsQuals["region"].GetStringValue()

	// Handle empty name or region
	if name == "" || regionName == "" {
		return nil, nil
	}

	// Create service
	svc, err := ELBV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_ssl_policy.getEc2SslPolicy", "connection_error", err)
		return nil, err
	}

	// Build params
	params := &elasticloadbalancingv2.DescribeSSLPoliciesInput{
		Names: []string{name},
	}

	if matrixKeyRegion == regionName {
		op, err := svc.DescribeSSLPolicies(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ec2_ssl_policy.getEc2SslPolicy", "api_error", err)
			return nil, err
		}

		if len(op.SslPolicies) > 0 {
			return op.SslPolicies[0], nil
		}
	}

	return nil, nil
}

func getEc2SslPolicyAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	data := h.Item.(types.SslPolicy)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":elbv2:" + region + ":" + commonColumnData.AccountId + ":ssl-policy/" + *data.Name}

	return akas, nil
}
