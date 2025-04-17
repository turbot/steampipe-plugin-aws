package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rolesanywhere"

	rolesanywherev1 "github.com/aws/aws-sdk-go/service/rolesanywhere"

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
			KeyColumns: plugin.AnyColumn([]string{"id"}),
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
		GetMatrixItemFunc: SupportedRegionMatrix(rolesanywherev1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The ID of the trust anchor.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TrustAnchorId"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the trust anchor.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TrustAnchorArn"),
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
				Description: "The trust anchor type.",
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
				Description: "Trust anchor expiry notification settings.",
				Type:        proto.ColumnType_JSON,
			},
		}),
	}
}

func getTrustAnchor(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	anchor_id := d.EqualsQuals["id"].GetStringValue()

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

			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
