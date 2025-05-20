package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rolesanywhere"
	"github.com/aws/aws-sdk-go-v2/service/rolesanywhere/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsRolesAnywhereTrustAnchor(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rolesanywhere_trust_anchor",
		Description: "AWS Roles Anywhere Trust Anchor",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"trust_anchor_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "ValidationException"}),
			},
			Hydrate: getTrustAnchor,
			Tags:    map[string]string{"service": "rolesanywhere", "action": "GetTrustAnchor"},
		},
		List: &plugin.ListConfig{
			Hydrate: listTrustAnchors,
			Tags:    map[string]string{"service": "rolesanywhere", "action": "ListTrustAnchors"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: listTagsForTrustAnchor,
				Tags: map[string]string{"service": "rolesanywhere", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_ROLESANYWHERE_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "trust_anchor_id",
				Description: "The unique identifier of the trust anchor.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the trust anchor.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TrustAnchorArn"),
			},
			{
				Name:        "name",
				Description: "The name of the trust anchor.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "updated_at",
				Description: "The date and time when the trust anchor was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "created_at",
				Description: "The date and time when the trust anchor was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "enabled",
				Description: "If the trust anchor is enabled or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "source_type",
				Description: "The type of the trust anchor.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Source.SourceType"),
			},
			{
				Name:        "source_data",
				Description: "The certificate/arn data for the trust anchor, depending on the source type.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Source.SourceData.Value"),
			},
			{
				Name:        "notification_settings",
				Description: "A list of notification settings to be associated to the trust anchor.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the anchor",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForTrustAnchor,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForTrustAnchor,
				Transform:   transform.From(trustAnchorTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("TrustAnchorArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listTrustAnchors(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	svc, err := RolesAnywhereClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rolesanywhere_trust_anchor.listTrustAnchors", "client_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}

	input := rolesanywhere.ListTrustAnchorsInput{}
	paginator := rolesanywhere.NewListTrustAnchorsPaginator(svc, &input, func(o *rolesanywhere.ListTrustAnchorsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_rolesanywhere_trust_anchor.listTrustAnchors", "api_error", err)
			return nil, err
		}

		for _, anchor := range output.TrustAnchors {
			d.StreamListItem(ctx, anchor)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getTrustAnchor(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	anchor_id := d.EqualsQuals["trust_anchor_id"].GetStringValue()

	svc, err := RolesAnywhereClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rolesanywhere_trust_anchor.getTrustAnchor", "client_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}

	params := &rolesanywhere.GetTrustAnchorInput{
		TrustAnchorId: aws.String(anchor_id),
	}

	op, err := svc.GetTrustAnchor(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rolesanywhere_trust_anchor.getTrustAnchor", "api_error", err)
		return nil, err
	}
	return *op.TrustAnchor, nil
}

func listTagsForTrustAnchor(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	anchor_arn := h.Item.(types.TrustAnchorDetail).TrustAnchorArn

	svc, err := RolesAnywhereClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rolesanywhere_trust_anchor.listTagsForTrustAnchor", "client_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}

	params := &rolesanywhere.ListTagsForResourceInput{
		ResourceArn: anchor_arn,
	}
	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rolesanywhere_trust_anchor.listTagsForTrustAnchor", "api_error", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTION

func trustAnchorTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tags := d.HydrateItem.(*rolesanywhere.ListTagsForResourceOutput)
	var turbotTagsMap map[string]string
	if tags.Tags != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tags.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}
	return turbotTagsMap, nil
}
