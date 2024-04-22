package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"

	cloudtrailv1 "github.com/aws/aws-sdk-go/service/cloudtrail"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudtrailImport(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudtrail_import",
		Description: "AWS CloudTrail Import",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("import_id"),
			Hydrate:    getCloudTrailImport,
			Tags:       map[string]string{"service": "cloudtrail", "action": "GetImport"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"UnsupportedOperationException", "ImportNotFoundException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudTrailImports,
			Tags:    map[string]string{"service": "cloudtrail", "action": "ListImports"},
			// For the location where the API operation is not supported, we receive UnsupportedOperationException.
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"UnsupportedOperationException"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "import_status",
					Require: plugin.Optional,
				},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getCloudTrailImport,
				Tags: map[string]string{"service": "cloudtrail", "action": "GetImport"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(cloudtrailv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "import_id",
				Description: "The ID of the import.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_timestamp",
				Description: "The timestamp of the import's creation.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "import_status",
				Description: "The status of the import.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "end_event_time",
				Description: "Used with EndEventTime to bound a StartImport request, and limit imported trail events to only those events logged within a specified time period.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getCloudTrailImport,
			},
			{
				Name:        "start_event_time",
				Description: "Used with StartEventTime to bound a StartImport request, and limit imported trail events to only those events logged within a specified time period.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getCloudTrailImport,
			},
			{
				Name:        "updated_timestamp",
				Description: "The timestamp of the import's last update.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "destinations",
				Description: "The ARN of the destination event data store.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "import_source",
				Description: "The source S3 bucket.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudTrailImport,
				Transform:   transform.FromField("ImportSource.S3"),
			},
			{
				Name:        "import_statistics",
				Description: "Provides statistics for the import.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudTrailImport,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImportId"),
			},
		}),
	}
}

//// LIST FUNCTION

func listCloudTrailImports(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := CloudTrailClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudtrail_import.listCloudTrailImports", "client_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &cloudtrail.ListImportsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	if d.EqualsQualString("import_status") != "" {
		input.ImportStatus = types.ImportStatus(d.EqualsQualString("import_status"))
	}

	paginator := cloudtrail.NewListImportsPaginator(svc, input, func(o *cloudtrail.ListImportsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudtrail_import.listCloudTrailImports", "api_error", err)
			return nil, err
		}

		for _, item := range output.Imports {
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCloudTrailImport(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var importId string
	if h.Item != nil {
		importId = *h.Item.(types.ImportsListItem).ImportId
	} else {
		importId = d.EqualsQuals["import_id"].GetStringValue()
	}

	// Empty check
	if importId == "" {
		return nil, nil
	}

	// Create session
	svc, err := CloudTrailClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudtrail_import.getCloudTrailImport", "client_error", err)
		return nil, err
	}

	params := &cloudtrail.GetImportInput{
		ImportId: &importId,
	}

	// execute list call
	op, err := svc.GetImport(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudtrail_import.getCloudTrailImport", "api_error", err)
		return nil, err
	}

	return op, nil
}
