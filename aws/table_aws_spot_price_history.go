package aws

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
	"time"
)

//// TABLE DEFINITION

func tableAwsSpotPriceHistory(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_spot_price_history",
		Description: "AWS Spot Price History",
		List: &plugin.ListConfig{
			Hydrate: listSpotPriceHistory,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "availability_zone", Require: plugin.Optional},
				{Name: "instance_type", Require: plugin.Optional},
				{Name: "product_description", Require: plugin.Optional},
				{Name: "start_time", Require: plugin.Optional},
				{Name: "end_time", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{Name: "availability_zone", Description: "The Availability Zone.", Type: proto.ColumnType_STRING},
			{Name: "instance_type", Description: "The instance type.", Type: proto.ColumnType_STRING},
			{Name: "product_description", Description: "A general description of the AMI.", Type: proto.ColumnType_STRING},
			{Name: "spot_price", Description: "The maximum price per unit hour that you are willing to pay for a Spot Instance.", Type: proto.ColumnType_STRING},
			{Name: "timestamp", Description: "The date and time the request was created", Type: proto.ColumnType_TIMESTAMP},
			{Name: "start_time", Description: "The date and time, up to the past 90 days, from which to start retrieving the price history data.", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromQual("start_time")},
			{Name: "end_time", Description: "The date and time, up to the current date, from which to stop retrieving the price history data.", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromQual("end_time")},
		}),
	}
}

//// LIST FUNCTION

func listSpotPriceHistory(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_spot_price_history.listSpotPriceHistory", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxItems := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 20 {
				maxItems = 20
			} else {
				maxItems = limit
			}
		}
	}

	input := ec2.DescribeSpotPriceHistoryInput{
		MaxResults: &maxItems,
	}

	// use filters instead of availability_zone which doesn't work
	if d.Quals["availability_zone"] != nil {
		value := getQualsValueByColumn(d.Quals, "availability_zone", "string").(string)
		columnName := "availability-zone"
		input.Filters = append(input.Filters, types.Filter{Name: &columnName, Values: []string{value}})
	}

	if d.Quals["instance_type"] != nil {
		value := getQualsValueByColumn(d.Quals, "instance_type", "string")
		v, ok := value.(string)
		if !ok {
			err := errors.New("instance_type must be a string")
			plugin.Logger(ctx).Error("aws_spot_price_history.listSpotPriceHistory", "input_type_error", err, "value", v)
			return nil, err
		}
		input.InstanceTypes = []types.InstanceType{types.InstanceType(v)}
	}

	if d.Quals["product_description"] != nil {
		value := getQualsValueByColumn(d.Quals, "product_description", "string")
		v, ok := value.(string)
		if !ok {
			err := errors.New("product_description must be a string")
			plugin.Logger(ctx).Error("aws_spot_price_history.listSpotPriceHistory", "input_type_error", err, "value", v)
			return nil, err
		}
		input.ProductDescriptions = []string{v}
	}

	if d.Quals["start_time"] != nil {
		value := getQualsValueByColumn(d.Quals, "start_time", "time")
		v, ok := value.(time.Time)
		if !ok {
			err := errors.New("start_time must have a time value")
			plugin.Logger(ctx).Error("aws_spot_price_history.listSpotPriceHistory", "input_type_error", err, "value", v)
			return nil, err
		}
		input.StartTime = &v
	}

	if d.Quals["end_time"] != nil {
		value := getQualsValueByColumn(d.Quals, "end_time", "time")
		v, ok := value.(time.Time)
		if !ok {
			err := errors.New("end_time must have a time value")
			plugin.Logger(ctx).Error("aws_spot_price_history.listSpotPriceHistory", "input_type_error", err, "value", v)
			return nil, err
		}
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
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_spot_price_history.listSpotPriceHistory", "api_error", err)
			return nil, err
		}

		for _, spotPrice := range output.SpotPriceHistory {
			d.StreamListItem(ctx, spotPrice)
			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
