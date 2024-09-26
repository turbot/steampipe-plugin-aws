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
			KeyColumns: plugin.AnyColumn([]string{"resource_arn", "id"}),
			Hydrate:    getAwsShieldProtection,
			Tags:       map[string]string{"service": "shield", "action": "DescribeProtection"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsShieldProtections,
			KeyColumns: plugin.OptionalColumns([]string{"name", "resource_arn", "resource_type"}),
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
			{
				Name:        "resource_type",
				Description: "The type of protected resource whose protections you want to retrieve.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("resource_type"),
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
		InclusionFilters: &types.InclusionProtectionFilters{},
	}

	if d.Quals["name"] != nil {
		for _, q := range d.Quals["name"].Quals {
			input.InclusionFilters.ProtectionNames = []string{}
			input.InclusionFilters.ProtectionNames = append(input.InclusionFilters.ProtectionNames, q.Value.GetStringValue())
		}
	}

	if d.Quals["resource_arn"] != nil {
		for _, q := range d.Quals["resource_arn"].Quals {
			input.InclusionFilters.ResourceArns = []string{}
			input.InclusionFilters.ResourceArns = append(input.InclusionFilters.ResourceArns, q.Value.GetStringValue())
		}
	}

	if d.Quals["resource_type"] != nil {
		for _, q := range d.Quals["resource_type"].Quals {
			input.InclusionFilters.ResourceTypes = []types.ProtectedResourceType{}
			input.InclusionFilters.ResourceTypes = append(input.InclusionFilters.ResourceTypes, types.ProtectedResourceType(q.Value.GetStringValue()))
		}
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
	var resourceArn string

	if h.Item != nil {
		protection := h.Item.(types.Protection)
		protectionId = *protection.Id
		resourceArn = *protection.ResourceArn
	} else {
		protectionId = d.EqualsQualString("id")
		resourceArn = d.EqualsQualString("resource_arn")
	}

	var params *shield.DescribeProtectionInput
	if protectionId != "" {
		params = &shield.DescribeProtectionInput{
			ProtectionId: aws.String(protectionId),
		}
	} else if resourceArn != "" {
		params = &shield.DescribeProtectionInput{
			ResourceArn: aws.String(resourceArn),
		}
	} else {
		return nil, nil
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
