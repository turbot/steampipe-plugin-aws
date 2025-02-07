package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ec2Endpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

//// TABLE DEFINITION

func tableAwsEc2SpotPrice(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_spot_price",
		Description: "AWS EC2 Spot Price History",
		List: &plugin.ListConfig{
			Hydrate: listEc2SpotPrice,
			Tags:    map[string]string{"service": "ec2", "action": "DescribeSpotPriceHistory"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValue"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "availability_zone", Require: plugin.Optional},
				{Name: "instance_type", Require: plugin.Optional},
				{Name: "product_description", Require: plugin.Optional},
				{Name: "start_time", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "end_time", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2Endpoint.AWS_EC2_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{Name: "availability_zone", Description: "The Availability Zone.", Type: proto.ColumnType_STRING},
			{Name: "instance_type", Description: "The instance type.", Type: proto.ColumnType_STRING},
			{Name: "product_description", Description: "A general description of the AMI.", Type: proto.ColumnType_STRING},
			{Name: "spot_price", Description: "The maximum price per unit hour that you are willing to pay for a Spot Instance.", Type: proto.ColumnType_STRING},
			{Name: "create_timestamp", Description: "The time stamp of the Spot price history.", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Timestamp")},
			{Name: "start_time", Description: "The date and time, up to the past 90 days, from which to start retrieving the price history data.", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromQual("start_time")},
			{Name: "end_time", Description: "The date and time, up to the current date, from which to stop retrieving the price history data.", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromQual("end_time")},
		}),
	}
}

//// LIST FUNCTION

func listEc2SpotPrice(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_spot_price_history.listEc2SpotPrice", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxItems := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = 1
			} else {
				maxItems = limit
			}
		}
	}

	input := ec2.DescribeSpotPriceHistoryInput{
		MaxResults: &maxItems,
	}

	equalQuals := d.EqualsQuals
	if d.Quals["availability_zone"] != nil {
		input.AvailabilityZone = aws.String(equalQuals["availability_zone"].GetStringValue())
	}

	if d.Quals["instance_type"] != nil {
		input.InstanceTypes = []types.InstanceType{types.InstanceType(equalQuals["instance_type"].GetStringValue())}
	}

	if d.Quals["product_description"] != nil {
		input.ProductDescriptions = []string{equalQuals["product_description"].GetStringValue()}
	}

	if d.Quals["start_time"] != nil {
		v := equalQuals["start_time"].GetTimestampValue().AsTime()
		input.StartTime = &v
	}

	if d.Quals["end_time"] != nil {
		v := equalQuals["end_time"].GetTimestampValue().AsTime()
		input.EndTime = &v
	}

	paginator := ec2.NewDescribeSpotPriceHistoryPaginator(
		svc,
		&input,
		func(o *ec2.DescribeSpotPriceHistoryPaginatorOptions) {
			o.Limit = maxItems
			o.StopOnDuplicateToken = true
		},
	)

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_spot_price_history.listEc2SpotPrice", "api_error", err)
			return nil, err
		}

		for _, spotPrice := range output.SpotPriceHistory {
			d.StreamListItem(ctx, spotPrice)
			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
