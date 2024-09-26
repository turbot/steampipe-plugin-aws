package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/shield"
	"github.com/aws/aws-sdk-go-v2/service/shield/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsShieldProtection(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_shield_protection",
		Description: "AWS Shield Protection",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getAwsShieldProtection,
			Tags:       map[string]string{"service": "shield", "action": "DescribeProtection"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsShieldProtections,
			Tags:    map[string]string{"service": "shield", "action": "ListProtections"},
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The unique identifier (ID) of the protection.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},
			{
				Name:        "name",
				Description: "The name of the protection.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "resource_arn",
				Description: "The ARN (Amazon Resource Name) of the Amazon Web Services resource that is protected.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceArn"),
			},
			{
				Name:        "health_check_ids",
				Description: "The unique identifier (ID) for the Route 53 health check that's associated with the protection.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HealthCheckIds").Transform(transform.EnsureStringArray),
			},
			{
				Name:        "arn",
				Description: "The ARN (Amazon Resource Name) of the protection.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProtectionArn"),
			},
			{
				Name:        "application_layer_automatic_response_configuration",
				Description: "The automatic application layer DDoS mitigation settings for the protection. This configuration determines whether Shield Advanced automatically manages rules in the web ACL in order to respond to application layer events that Shield Advanced determines to be DDoS attacks.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ApplicationLayerAutomaticResponseConfiguration"),
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
				Transform:   transform.FromField("ProtectionArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

func listAwsShieldProtections(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := ShieldClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_shield_protection.listAwsShieldProtections", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	queryResultLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		queryResultLimit = min(queryResultLimit, int32(*d.QueryContext.Limit))
	}

	input := &shield.ListProtectionsInput{
		MaxResults: aws.Int32(queryResultLimit),
	}
	paginator := shield.NewListProtectionsPaginator(svc, input, func(o *shield.ListProtectionsPaginatorOptions) {
		o.Limit = queryResultLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_shield_protection.listAwsShieldProtections", "api_error", err)
			return nil, err
		}

		for _, items := range output.Protections {
			d.StreamListItem(ctx, items)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getAwsShieldProtection(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := ShieldClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_shield_protection.getAwsShieldProtection", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	var protectionId string
	if h.Item != nil {
		protection := h.Item.(types.Protection)
		protectionId = *protection.Id
	} else {
		protectionId = d.EqualsQualString("id")
	}

	params := &shield.DescribeProtectionInput{
		ProtectionId: aws.String(protectionId),
	}

	data, err := svc.DescribeProtection(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_shield_protection.getAwsProtection", "api_error", err)
		return nil, err
	}

	if data != nil {
		return data.Protection, nil
	}

	return nil, nil
}
