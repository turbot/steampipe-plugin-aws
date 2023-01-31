package aws

import (
	"context"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mgn"
	"github.com/aws/aws-sdk-go-v2/service/mgn/types"

	mgnv1 "github.com/aws/aws-sdk-go/service/mgn"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINATION

func tableAwsMGNApplication(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_mgn_application",
		Description: "AWS MGN Application",
		List: &plugin.ListConfig{
			Hydrate: ListAwsMGNApplications,
			IgnoreConfig: &plugin.IgnoreConfig{
				// UninitializedAccountException - This error comes up when the service is not enabled for the account.
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"UninitializedAccountException"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "application_id", Require: plugin.Optional},
				{Name: "wave_id", Require: plugin.Optional},
				{Name: "is_archived", Require: plugin.Optional, Operators: []string{"=", "<>"}},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(mgnv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "Application name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "application_id",
				Description: "Application ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationID"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_date_time",
				Description: "Application creation dateTime.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "Application description.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_archived",
				Description: "Application archival status.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "last_modified_date_time",
				Description: "Application last modified dateTime.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "wave_id",
				Description: "Application wave ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WaveID"),
			},
			{
				Name:        "application_aggregated_status",
				Description: "Application aggregated status.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: "A list of tags attached to the application.",
				Type:        proto.ColumnType_JSON,
			},
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
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

// // LIST FUNCTION
func ListAwsMGNApplications(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	svc, err := MGNClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_mgn_application.ListAwsMGNApplications", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	// Page size must be greater than 0 and less than or equal to 1000
	input := &mgn.ListApplicationsInput{
		MaxResults: maxLimit,
	}

	filter := buildMGNApplicationFilter(d.Quals)
	if filter != nil {
		input.Filters = filter
	}

	paginator := mgn.NewListApplicationsPaginator(svc, input, func(o *mgn.ListApplicationsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_mgn_application.ListAwsMGNApplications", "api_error", err)
			return nil, err
		}

		for _, item := range output.Items {
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// UTILITY FUNCTION

// Build MGN Application list call input filter
func buildMGNApplicationFilter(quals plugin.KeyColumnQualMap) *types.ListApplicationsRequestFilters {
	columnNames := []string{"application_id", "is_archived", "wave_id"}

	filter := &types.ListApplicationsRequestFilters{}
	for _, columnName := range columnNames {
		if quals[columnName] != nil {
			switch columnName {
			case "application_id":
				value := getQualsValueByColumn(quals, columnName, "string")
				if value != nil {
					val, ok := value.(string)
					if ok {
						filter.ApplicationIDs = []string{val}
					}
				}
			case "wave_id":
				value := getQualsValueByColumn(quals, columnName, "string")
				if value != nil {
					val, ok := value.(string)
					if ok {
						filter.WaveIDs = []string{val}
					}
				}
			case "is_archived":
				value := getQualsValueByColumn(quals, columnName, "boolean")
				boolValue, err := strconv.ParseBool(value.(string))
				if err != nil {
					log.Fatal(err)
				}
				filter.IsArchived = aws.Bool(boolValue)
			}
		}
	}

	return filter
}
