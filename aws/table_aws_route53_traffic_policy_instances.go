package aws

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
)

func tableAwsRoute53TrafficPolicyInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_route53_traffic_policy_instances",
		Description: "AWS Route53 Traffic Policy Instances",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"id", "version"}),
			Hydrate:           getTrafficPolicyInstance,
			ShouldIgnoreError: isNotFoundError([]string{"NoSuchTrafficPolicyInstance"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listTrafficPolicyInstances,
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name that you specified when traffic policy was created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The ID that Amazon Route 53 assigned to a traffic policy when it was created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The DNS type of the resource record sets that Amazon Route 53 creates when you use a traffic policy to create a traffic policy instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version",
				Description: "The version number that Amazon Route 53 assigns to a traffic policy.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "comment",
				Description: "The comment that you specified when traffic policy was created.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTrafficPolicyInstance,
			},
			{
				Name:        "document",
				Description: "The definition of a traffic policy in JSON format.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTrafficPolicyInstance,
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
				Hydrate:     getRoute53TrafficPolicyInstanceTurbotAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listTrafficPolicyInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listTrafficPolicyInstances")

	// Create session
	svc, err := Route53Service(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &route53.ListTrafficPolicyInstancesInput{
		MaxItems: aws.String("100"),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < 100 {
			if *limit < 1 {
				input.MaxItems = aws.String("1")
			} else {
				input.MaxItems = aws.String(fmt.Sprint(*limit))
			}
		}
	}

	// List call
	pagesLeft := true
	for pagesLeft {
		result, err := svc.ListTrafficPolicyInstances(input)
		if err != nil {
			plugin.Logger(ctx).Error("listTrafficPolicyInstances", "ListTrafficPolicyInstances_error", err)
			return nil, err
		}

		for _, policies := range result.TrafficPolicyInstances {
			d.StreamListItem(ctx, policies)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				pagesLeft = false
			}
		}

		// wait for all executions to be processed
		if *result.IsTruncated {
			input.TrafficPolicyInstanceNameMarker = result.TrafficPolicyInstanceNameMarker
		} else {
			pagesLeft = false
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getTrafficPolicyInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getTrafficPolicyInstance")

	instance := h.Item.(*route53.TrafficPolicyInstance)
	var id string
	if h.Item != nil {
		id = *instance.Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
	}

	// Validate if input params are empty
	if len(id) < 1 {
		return nil, nil
	}

	// Create session
	svc, err := Route53Service(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &route53.GetTrafficPolicyInstanceInput{
		Id: aws.String(id),
	}

	// execute get call
	item, err := svc.GetTrafficPolicyInstance(params)
	if err != nil {
		return nil, err
	}
	return item.TrafficPolicyInstance, nil
}

func getRoute53TrafficPolicyInstanceTurbotAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instance := h.Item.(*route53.TrafficPolicyInstance)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	//arn:aws:route53::<account-id>:trafficpolicy/<id>/<version>
	akas := []string{"arn:" + commonColumnData.Partition +
		":route53::" + commonColumnData.AccountId +
		":" + "trafficpolicy/" + *instance.Id +
		"/" + strconv.FormatInt(trafficPolicyVersion(h.Item), 10)}

	return akas, nil
}
