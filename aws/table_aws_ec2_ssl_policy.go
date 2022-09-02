package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsEc2SslPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_ssl_policy",
		Description: "AWS EC2 SSL Policy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "region"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"SSLPolicyNotFound"}),
			},
			Hydrate: getEc2SslPolicy,
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2SslPolicies,
		},
		GetMatrixItemFunc: BuildRegionList,
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
	plugin.Logger(ctx).Trace("listEc2SslPolicies")

	// Create Session
	svc, err := ELBv2Service(ctx, d)
	if err != nil {
		return nil, err
	}

	// List call
	params := &elbv2.DescribeSSLPoliciesInput{
		PageSize: aws.Int64(400),
	}

	// Limiting the results
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *params.PageSize {
			if *limit < 1 {
				params.PageSize = aws.Int64(1)
			} else {
				params.PageSize = limit
			}
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := svc.DescribeSSLPolicies(params)
		if err != nil {
			return nil, err
		}

		for _, sslPolicy := range response.SslPolicies {
			d.StreamListItem(ctx, sslPolicy)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
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
	plugin.Logger(ctx).Trace("getEc2SslPolicy")

	matrixKeyRegion := d.KeyColumnQualString(matrixKeyRegion)
	name := d.KeyColumnQuals["name"].GetStringValue()
	regionName := d.KeyColumnQuals["region"].GetStringValue()

	// Handle empty name or region
	if name == "" || regionName == "" {
		return nil, nil
	}

	// Create service
	svc, err := ELBv2Service(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build params
	params := &elbv2.DescribeSSLPoliciesInput{
		Names: []*string{aws.String(name)},
	}

	if matrixKeyRegion == regionName {
		op, err := svc.DescribeSSLPolicies(params)
		if err != nil {
			return nil, err
		}

		if op.SslPolicies != nil && len(op.SslPolicies) > 0 {
			return op.SslPolicies[0], nil
		}
	}

	return nil, nil
}

func getEc2SslPolicyAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEc2SslPolicyAkas")
	region := d.KeyColumnQualString(matrixKeyRegion)
	data := h.Item.(*elbv2.SslPolicy)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":elbv2:" + region + ":" + commonColumnData.AccountId + ":ssl-policy/" + *data.Name}

	return akas, nil
}
