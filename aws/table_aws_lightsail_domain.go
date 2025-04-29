package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/aws/aws-sdk-go-v2/service/lightsail/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsLightsailDomain(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_lightsail_domain",
		Description: "AWS Lightsail Domain",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("domain_name"),
			Hydrate:    getLightsailDomain,
			Tags:       map[string]string{"service": "lightsail", "action": "GetDomain"},
		},
		List: &plugin.ListConfig{
			Hydrate: listLightsailDomains,
			Tags:    map[string]string{"service": "lightsail", "action": "GetDomains"},
		},
		GetMatrixItemFunc: func(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {
			// Lightsail domain API is only available in us-east-1
			matrix := []map[string]interface{}{}
			commonData, err := getCommonColumns(ctx, d, nil)
			if err != nil {
				return matrix
			}
			commonColumnData := commonData.(*awsCommonColumnData)
			matrix = append(matrix, map[string]interface{}{
				"region":    "us-east-1",
				"account":   commonColumnData.AccountId,
				"partition": commonColumnData.Partition,
			})
			return matrix
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "domain_name",
				Description: "The name of the domain.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the domain.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Arn"),
			},
			{
				Name:        "created_at",
				Description: "The timestamp when the domain was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreatedAt"),
			},
			{
				Name:        "is_managed",
				Description: "Indicates whether the domain is managed by Lightsail.",
				Type:        proto.ColumnType_BOOL,
				Transform: transform.FromField("RegisteredDomainDelegationInfo.NameServersUpdateState.Code").Transform(transform.ToString).Transform(func(ctx context.Context, d *transform.TransformData) (interface{}, error) {
					if d.Value == nil {
						return false, nil
					}
					return d.Value.(string) == "SUCCEEDED", nil
				}),
			},
			{
				Name:        "location",
				Description: "The AWS Region and Availability Zones where the domain is located.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Location"),
			},
			{
				Name:        "resource_type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceType"),
			},
			{
				Name:        "support_code",
				Description: "The support code. Include this code in your email to support when you have questions about your Lightsail domain.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SupportCode"),
			},
			{
				Name:        "domain_entries",
				Description: "An array of key-value pairs containing information about the domain entries.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DomainEntries"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the domain.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getLightsailDomainTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listLightsailDomains(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := LightsailDomainClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lightsail_domain.listLightsailDomains", "connection_error", err)
		return nil, err
	}

	input := &lightsail.GetDomainsInput{}

	// List call
	for {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := svc.GetDomains(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_lightsail_domain.listLightsailDomains", "api_error", err)
			return nil, err
		}

		plugin.Logger(ctx).Debug("aws_lightsail_domain.listLightsailDomains", "domains_count", len(output.Domains))
		for _, domain := range output.Domains {
			plugin.Logger(ctx).Debug("aws_lightsail_domain.listLightsailDomains", "domain", domain)
			d.StreamListItem(ctx, domain)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if output.NextPageToken == nil {
			break
		}
		input.PageToken = output.NextPageToken
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getLightsailDomain(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := LightsailDomainClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lightsail_domain.getLightsailDomain", "connection_error", err)
		return nil, err
	}

	name := d.EqualsQuals["domain_name"].GetStringValue()

	if name == "" {
		return nil, nil
	}

	params := &lightsail.GetDomainInput{
		DomainName: aws.String(name),
	}

	op, err := svc.GetDomain(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Debug("aws_lightsail_domain.getLightsailDomain", "api_error", err)
		return nil, err
	}

	if op.Domain != nil {
		return op.Domain, nil
	}

	return nil, nil
}

func getLightsailDomainTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.(types.Domain).Tags

	if len(tags) == 0 {
		return nil, nil
	}

	turbotTagsMap := map[string]string{}
	for _, i := range tags {
		turbotTagsMap[*i.Key] = *i.Value
	}

	return turbotTagsMap, nil
}

// // ////////////////////////////////////////////////////////////////////////////
// // TRANSFORM FUNCTIONS
// // ////////////////////////////////////////////////////////////////////////////

func domainNameToARN(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}
	domainName, ok := d.Value.(string)
	if !ok {
		return nil, nil
	}
	accountId, ok := d.MatrixItem["account"].(string)
	if !ok || accountId == "" {
		return nil, nil
	}
	return "arn:aws:lightsail:us-east-1:" + accountId + ":Domain/" + domainName, nil
}
