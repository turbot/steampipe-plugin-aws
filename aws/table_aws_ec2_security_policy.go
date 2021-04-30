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

func tableAwsEc2SecurityPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_security_policy",
		Description: "AWS EC2 Security Policy",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("arn"),
			ShouldIgnoreError: isNotFoundError([]string{"name"}),
			Hydrate:           getEc2SecurityPolicy,
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2SecurityPolicies,
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
				Hydrate:     getAwsEc2SecurityPolicyAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listEc2SecurityPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listEc2SecurityPolicies", "AWS_REGION", region)

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

func getEc2SecurityPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEc2SecurityPolicy")

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

func getAwsEc2SecurityPolicyAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsEc2SecurityPolicyAkas")
	data := h.Item.(*elbv2.SslPolicy)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":elbv2:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":ssl-policy/" + *data.Name}

	return akas, nil
}
