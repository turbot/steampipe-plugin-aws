package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/oam"
	"github.com/aws/aws-sdk-go-v2/service/oam/types"

	oamv1 "github.com/aws/aws-sdk-go/service/oam"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsOAMLink(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_oam_link",
		Description: "AWS OAM Link",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterException", "ResourceNotFoundException"}),
			},
			Hydrate: getAwsOAMLink,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsOAMLinks,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(oamv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The random ID string that Amazon Web Service generates as part of the link ARN.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The ARN of the link.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sink_arn",
				Description: "The ARN of the sink that this link is attached to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "label",
				Description: "The label that was assigned to this link at creation, with the variables resolved to their actual values.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "label_template",
				Description: "The exact label template that was specified when the link was created, with the template variables not resolved.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsOAMLink,
			},
			{
				Name:        "resource_types",
				Description: "The resource types supported by this link.",
				Type:        proto.ColumnType_JSON,
			},

			/// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsOAMLink,
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

func listAwsOAMLinks(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create client
	svc, err := OAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_oam_link.listAwsOAMLinks", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(5)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &oam.ListLinksInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := oam.NewListLinksPaginator(svc, input, func(o *oam.ListLinksPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_oam_link.listAwsOAMLinks", "api_error", err)
			return nil, err
		}

		for _, link := range output.Items {
			d.StreamListItem(ctx, link)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsOAMLink(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var arn string
	if h.Item != nil {
		arn = *h.Item.(types.ListLinksItem).Arn
	} else {
		arn = d.EqualsQualString("arn")
	}

	// Empty Check
	if arn == "" {
		return nil, nil
	}

	// Create Client
	svc, err := OAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_oam_link.getAwsOAMLink", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &oam.GetLinkInput{
		Identifier: aws.String(arn),
	}

	// Get call
	resp, err := svc.GetLink(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_oam_link.getAwsOAMLink", "api_error", err)
		return nil, err
	}

	return resp, nil
}
