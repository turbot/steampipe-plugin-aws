package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"

	rdsv1 "github.com/aws/aws-sdk-go/service/rds"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsRDSDBRecommendation(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_db_recommendation",
		Description: "AWS RDS DB Recommendation",
		List: &plugin.ListConfig{
			Hydrate: listRDSDBRecommendations,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "recommendation_id", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
				{Name: "severity", Require: plugin.Optional},
				{Name: "type_id", Require: plugin.Optional},
				{Name: "updated_time", Require: plugin.Optional, Operators: []string{">=", "<="}},
			},
			Tags: map[string]string{"service": "rds", "action": "DescribeDBRecommendations"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(rdsv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "recommendation_id",
				Description: "The unique identifier for the recommendation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_arn",
				Description: "The Amazon Resource Name (ARN) of the RDS resource associated with the recommendation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceArn"),
			},
			{
				Name:        "category",
				Description: "The category of the recommendation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The current status of the recommendation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_time",
				Description: "The time when the recommendation was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated_time",
				Description: "The time when the recommendation was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("UpdatedTime").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "description",
				Description: "A detailed description of the recommendation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "detection",
				Description: "A short description of the issue identified for this recommendation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "impact",
				Description: "A short description that explains the possible impact of an issue.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "reason",
				Description: "The reason why this recommendation was created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "recommendation",
				Description: "A short description of the recommendation to resolve an issue.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "severity",
				Description: "The severity level of the recommendation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source",
				Description: "The AWS service that generated the recommendations.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type_detection",
				Description: "A short description of the recommendation type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type_id",
				Description: "A value that indicates the type of recommendation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type_recommendation",
				Description: "A short description that summarizes the recommendation to fix all the issues of the recommendation type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "additional_info",
				Description: "Additional information about the recommendation.",
				Type:        proto.ColumnType_STRING,
			},

			// JSON fields
			{
				Name:        "issue_details",
				Description: "Details of the issue that caused the recommendation.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "links",
				Description: "A link to documentation that provides additional information about the recommendation.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "recommended_actions",
				Description: "A list of recommended actions.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RecommendationId"),
			},
		}),
	}
}

//// LIST FUNCTION

func listRDSDBRecommendations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := RDSDBRecommendationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_recommendation.listRDSDBRecommendations", "connection_error", err)
		return nil, err
	}
	// Unsupported region check
	if svc == nil {
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(50)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &rds.DescribeDBRecommendationsInput{
		MaxRecords: aws.Int32(maxLimit),
	}

	if d.Quals["updated_time"] != nil {
		for _, q := range d.Quals["updated_time"].Quals {
			value := q.Value.GetTimestampValue().AsTime()
			if q.Operator == ">=" {
				input.LastUpdatedAfter = &value
			}
			if q.Operator == "<=" {
				input.LastUpdatedBefore = &value
			}
		}
	}

	// Build input  filter parameter
	filters := buildRdsDBRecommendationFilter(d.Quals)

	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := rds.NewDescribeDBRecommendationsPaginator(svc, input, func(o *rds.DescribeDBRecommendationsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_rds_db_recommendation.listRDSDBRecommendations", "api_error", err)
			return nil, err
		}

		for _, recommendation := range output.DBRecommendations {
			d.StreamListItem(ctx, recommendation)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// UTILITY FUNCTIONS

// Build RDS DB Recommendation list call input filter
func buildRdsDBRecommendationFilter(quals plugin.KeyColumnQualMap) []types.Filter {
	filters := make([]types.Filter, 0)
	// We have intentionally avoided using the other filter parameters:
	// - dbi-resource-id
	// - cluster-resource-id
	// - pg-arn
	// - cluster-pg-arn
	// because they are not part of the API response.
	filterQuals := map[string]string{
		"recommendation_id": "recommendation-id",
		"status":            "status",
		"severity":          "severity",
		"type_id":           "type-id",
	}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := types.Filter{
				Name: aws.String(filterName),
			}
			value := getQualsValueByColumn(quals, columnName, "string")
			val, ok := value.(string)
			if ok {
				filter.Values = []string{val}
			}
			filters = append(filters, filter)
		}
	}
	return filters
}
