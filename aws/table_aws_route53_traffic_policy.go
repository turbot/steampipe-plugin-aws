package aws

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsRoute53TrafficPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_route53_traffic_policy",
		Description: "AWS Route53 Traffic Policy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id", "version"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchTrafficPolicy"}),
			},
			Hydrate: getTrafficPolicy,
		},
		List: &plugin.ListConfig{
			Hydrate: listTrafficPolicies,
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
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
			},
			{
				Name:        "comment",
				Description: "The comment that you specified when traffic policy was created.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTrafficPolicy,
			},
			{
				Name:        "document",
				Description: "The definition of a traffic policy in JSON format.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTrafficPolicy,
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
				Hydrate:     getRoute53TrafficPolicyTurbotAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listTrafficPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := Route53Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_traffic_policy.listTrafficPolicies", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(100)
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

	input := &route53.ListTrafficPoliciesInput{
		MaxItems: aws.Int32(maxLimit),
	}

	// List call
	pagesLeft := true
	for pagesLeft {
		result, err := svc.ListTrafficPolicies(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_route53_traffic_policy.listTrafficPolicies", "api_error", err)
			return nil, err
		}

		var wg sync.WaitGroup
		errorCh := make(chan error, len(result.TrafficPolicySummaries))

		for _, policies := range result.TrafficPolicySummaries {
			wg.Add(1)
			go listTrafficPolicyVersionsAsync(ctx, d, svc, policies.Id, &wg, errorCh)
		}

		// wait for all executions to be processed
		wg.Wait()
		close(errorCh)

		for err := range errorCh {
			plugin.Logger(ctx).Error("aws_route53_traffic_policy.listTrafficPolicies", "listTrafficPolicyVersionsAsync_error", err)
			return nil, err
		}

		if result.IsTruncated {
			input.TrafficPolicyIdMarker = result.TrafficPolicyIdMarker
		} else {
			pagesLeft = false
		}
	}
	return nil, nil
}

// To fetch all available versions for a traffic policy
func listTrafficPolicyVersionsAsync(ctx context.Context, d *plugin.QueryData, svc *route53.Client, id *string, wg *sync.WaitGroup, errorCh chan error) {

	// Limiting the results
	maxLimit := int32(100)
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

	input := &route53.ListTrafficPolicyVersionsInput{
		Id:       id,
		MaxItems: aws.Int32(maxLimit),
	}

	defer wg.Done()

	// List call
	pagesLeft := true
	for pagesLeft {
		result, err := svc.ListTrafficPolicyVersions(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_route53_traffic_policy.listTrafficPolicyVersionsAsync", "ListTrafficPolicyVersions_api_error", err)
			errorCh <- err
		}
		for _, policies := range result.TrafficPolicies {
			d.StreamListItem(ctx, policies)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				pagesLeft = false
			}
		}
		if result.IsTruncated {
			input.TrafficPolicyVersionMarker = result.TrafficPolicyVersionMarker
		} else {
			pagesLeft = false
		}
	}
}

//// HYDRATE FUNCTIONS

func getTrafficPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var id string
	var version int32
	if h.Item != nil {
		trafficPolicy := h.Item.(types.TrafficPolicy)
		id = *trafficPolicy.Id
		version = *trafficPolicy.Version
	} else {
		id = d.EqualsQuals["id"].GetStringValue()
		version = int32(d.EqualsQuals["version"].GetInt64Value())
	}

	// Validate if input params are empty
	if len(id) < 1 || version < 1 {
		return nil, nil
	}

	// Create session
	svc, err := Route53Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_traffic_policy.getTrafficPolicy", "connection_error", err)
		return nil, err
	}

	params := &route53.GetTrafficPolicyInput{
		Id:      aws.String(id),
		Version: &version,
	}

	// execute get call
	item, err := svc.GetTrafficPolicy(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_traffic_policy.getTrafficPolicy", "api_error", err)
		return nil, err
	}
	return *item.TrafficPolicy, nil
}

func getRoute53TrafficPolicyTurbotAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	trafficPolicy := h.Item.(types.TrafficPolicy)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	//arn:aws:route53::<account-id>:trafficpolicy/<id>/<version>
	arn := fmt.Sprintf("arn:%s:route53::%s:trafficpolicy/%s/%s", commonColumnData.Partition, commonColumnData.AccountId, *trafficPolicy.Id, fmt.Sprint(*trafficPolicy.Version))

	return []string{arn}, nil
}
