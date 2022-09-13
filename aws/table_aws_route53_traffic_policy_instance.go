package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
)

func tableAwsRoute53TrafficPolicyInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_route53_traffic_policy_instance",
		Description: "AWS Route53 Traffic Policy Instance",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id"}),
			Hydrate:    getTrafficPolicyInstance,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"NoSuchTrafficPolicyInstance"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listTrafficPolicyInstances,
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The DNS name for which Amazon Route 53 responds to queries.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The id that Amazon Route 53 assigned to the new traffic policy instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "hosted_zone_id",
				Description: "The id of the hosted zone that Amazon Route 53 created resource record sets in.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "message",
				Description: "If State is Failed, an explanation of the reason for the failure.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "Current state of the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "traffic_policy_id",
				Description: "The ID of the traffic policy that Amazon Route 53 used to create resource record sets in the specified hosted zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "traffic_policy_type",
				Description: "The DNS type that Amazon Route 53 assigned to all of the resource record sets that it created for this traffic policy instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "traffic_policy_version",
				Description: "The version of the traffic policy that Amazon Route 53 used to create resource record sets in the specified hosted zone.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "ttl",
				Description: "The TTL that Amazon Route 53 assigned to all of the resource record sets that it created in the specified hosted zone.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("TTL"),
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
	// Create session
	svc, err := Route53Service(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_traffic_policy_instance.listTrafficPolicyInstances", "service_creation_error", err)
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
			plugin.Logger(ctx).Error("aws_route53_traffic_policy_instance.listTrafficPolicyInstances", "api_err", err)
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
	var id string
	if h.Item != nil {
		id = *h.Item.(*route53.TrafficPolicyInstance).Id
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
		plugin.Logger(ctx).Error("aws_route53_traffic_policy_instance.getTrafficPolicyInstance", "service_creation_error", err)
		return nil, err
	}

	params := &route53.GetTrafficPolicyInstanceInput{
		Id: aws.String(id),
	}

	// execute get call
	item, err := svc.GetTrafficPolicyInstance(params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_traffic_policy_instance.getTrafficPolicyInstance", "api_error", err)
		return nil, err
	}
	return item.TrafficPolicyInstance, nil
}

func getRoute53TrafficPolicyInstanceTurbotAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instance := h.Item.(*route53.TrafficPolicyInstance)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_traffic_policy_instance.getRoute53TrafficPolicyInstanceTurbotAkas", "api_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	//arn:aws:route53::<account-id>:trafficpolicyinstance/<id>
	akas := []string{"arn:" + commonColumnData.Partition +
		":route53::" + commonColumnData.AccountId +
		":" + "trafficpolicyinstance/" + *instance.Id}

	return akas, nil
}
