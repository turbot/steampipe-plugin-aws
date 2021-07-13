package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsEc2SslPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_ssl_policy",
		Description: "AWS EC2 SSL Policy",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			ShouldIgnoreError: isNotFoundError([]string{"name"}),
			Hydrate:           getEc2SslPolicy,
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2SslPolicies,
		},
		GetMatrixItem: BuildRegionList,
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

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listEc2SslPolicies", "AWS_REGION", region)

	// Create Session
	svc, err := ELBv2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	params := &elbv2.DescribeSSLPoliciesInput{}
	pagesLeft := true
	for pagesLeft {
		response, err := svc.DescribeSSLPolicies(params)
		if err != nil {
			return nil, err
		}

		for _, sslPolicy := range response.SslPolicies {
			d.StreamListItem(ctx, sslPolicy)
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

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	name := d.KeyColumnQuals["name"].GetStringValue()

	// Create service
	svc, err := ELBv2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build params
	params := &elbv2.DescribeSSLPoliciesInput{
		Names: []*string{aws.String(name)},
	}

	op, err := svc.DescribeSSLPolicies(params)
	if err != nil {
		return nil, err
	}

	if op.SslPolicies != nil && len(op.SslPolicies) > 0 {
		return op.SslPolicies[0], nil
	}

	return nil, nil
}

func getEc2SslPolicyAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEc2SslPolicyAkas")
	data := h.Item.(*elbv2.SslPolicy)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":elbv2:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":ssl-policy/" + *data.Name}

	return akas, nil
}
