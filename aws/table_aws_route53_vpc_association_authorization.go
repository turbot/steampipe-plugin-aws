package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsRoute53VPCAssociationAuthorization(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_route53_vpc_association_authorization",
		Description: "AWS Route53 VPC Association Authorization",
		List: &plugin.ListConfig{
			ParentHydrate: listHostedZones,
			Hydrate:       listVPCAssociationAuthorization,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchHostedZone"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "hosted_zone_id",
					Require: plugin.Optional,
				},
			},
			Tags: map[string]string{"service": "route53", "action": "ListVPCAssociationAuthorizations"},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "hosted_zone_id",
				Description: "The ID of the hosted zone for which you want a list of VPCs that can be associated with the hosted zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "(Private hosted zones only) The ID of an Amazon VPC.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VPCId"),
			},
			{
				Name:        "vpc_region",
				Description: "(Private hosted zones only) The region that an Amazon VPC was created in.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VPCRegion"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(authorizationTitle).Transform(transform.ToString),
			},
		}),
	}
}

type VPCAssociationAuthorizationResult struct {
	HostedZoneId *string
	VPCId        *string
	VPCRegion    *string
}

//// LIST FUNCTION
func listVPCAssociationAuthorization(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := Route53Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_route53_zone_vpc_association_authorization.listVPCAssociationAuthorization", "client_error", err)
		return nil, err
	}

	var hostedZoneId = d.EqualsQualString("hosted_zone_id")

	// Try to use ParentHydrate if no explicit value in quals
	if hostedZoneId == "" {
		zone := h.Item.(HostedZoneResult)
		if zone.Id == nil {
			return nil, nil
		}
		hostedZoneId = strings.Split(*zone.Id, "/")[2]
	}

	// https://docs.aws.amazon.com/Route53/latest/APIReference/API_ListVPCAssociationAuthorizations.html#API_ListVPCAssociationAuthorizations_RequestSyntax
	// Max 50 per page as per aws api docs
	maxItems := int32(50)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = int32(limit)
		}
	}

	input := &route53.ListVPCAssociationAuthorizationsInput{
		HostedZoneId: aws.String(hostedZoneId),
		MaxResults:   aws.Int32(maxItems),
	}

	for {
		result, err := svc.ListVPCAssociationAuthorizations(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_route53_zone_vpc_association_authorization.listVPCAssociationAuthorization", "api_error", err)
			return nil, err
		}

		for _, vpc := range result.VPCs {
			vpcRegion := string(vpc.VPCRegion)
			d.StreamListItem(ctx, &VPCAssociationAuthorizationResult{
				HostedZoneId: &hostedZoneId,
				VPCId:        vpc.VPCId,
				VPCRegion:    &vpcRegion,
			})

			// context cancelled
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if result.NextToken == nil {
			break
		}
		input.NextToken = result.NextToken
	}
	return nil, nil
}

// TRANSFORM FUNCTIONS
func authorizationTitle(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	result := d.HydrateItem.(*VPCAssociationAuthorizationResult)

	// same format as terraform builtin IDs "HOSTED_ZONE_ID:VPC_ID"
	// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_vpc_association_authorization#import
	akas := fmt.Sprintf("%s:%s", *result.HostedZoneId, *result.VPCId)
	return akas, nil
}
