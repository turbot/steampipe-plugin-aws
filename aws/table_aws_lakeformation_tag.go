package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lakeformation"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsLakeformationTag(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_lakeformation_tag",
		Description: "AWS Lake Formation Tag",
		Get: &plugin.GetConfig{
			Hydrate:    getLakeformationTag,
			KeyColumns: plugin.AllColumns([]string{"catalog_id", "tag_key"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"EntityNotFoundException"}),
			},
			Tags: map[string]string{"service": "lakeformation", "action": "GetLFTag"},
		},
		List: &plugin.ListConfig{
			Hydrate: listLakeformationTags,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "catalog_id", Require: plugin.Optional},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"EntityNotFoundException"}),
			},
			Tags: map[string]string{"service": "lakeformation", "action": "ListLFTags"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_LAKEFORMATION_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			// Principal Information
			{
				Name:        "catalog_id",
				Description: "The identifier for the Data Catalog. By default, the account ID. The Data Catalog is the persistent metadata store.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tag_key",
				Description: "The key-name for the LF-tag.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tag_values",
				Description: "A list of possible values an attribute can take.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TagKey"),
			},
		}),
	}
}

//// LIST FUNCTION

func listLakeformationTags(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := LakeFormationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lakeformation_tag.listLakeformationTags", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	maxItems := int32(100)
	input := &lakeformation.ListLFTagsInput{}

	// Reduce the request limit based on user input
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
		}
	}
	input.MaxResults = aws.Int32(maxItems)

	if d.EqualsQualString("catalog_id") != "" {
		cId := d.EqualsQualString("catalog_id")
		input.CatalogId = &cId
	}

	paginator := lakeformation.NewListLFTagsPaginator(svc, input, func(o *lakeformation.ListLFTagsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// Apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_lakeformation_tag.listLakeformationTags", "api_error", err)
			return nil, err
		}

		for _, tag := range output.LFTags {
			d.StreamListItem(ctx, tag)

			// Stop if the context is cancelled
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTION

func getLakeformationTag(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	catalogId := d.EqualsQualString("catalog_id")
	tagKey := d.EqualsQualString("tag_key")

	if tagKey == "" || catalogId == "" {
		return nil, nil
	}

	// Create session
	svc, err := LakeFormationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lakeformation_tag.getLakeformationTag", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	input := &lakeformation.GetLFTagInput{
		CatalogId: &catalogId,
		TagKey:    &tagKey,
	}

	res, err := svc.GetLFTag(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lakeformation_tag.getLakeformationTag", "api_error", err)
		return nil, err
	}

	return res, nil
}
