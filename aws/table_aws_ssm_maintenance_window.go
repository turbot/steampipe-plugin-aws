package aws

import (
	"context"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"

	ssmv1 "github.com/aws/aws-sdk-go/service/ssm"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsSSMMaintenanceWindow(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssm_maintenance_window",
		Description: "AWS SSM Maintenance Window",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("window_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"DoesNotExistException"}),
			},
			Hydrate: getAwsSSMMaintenanceWindow,
			Tags:    map[string]string{"service": "ssm", "action": "GetMaintenanceWindow"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsSSMMaintenanceWindow,
			Tags:    map[string]string{"service": "ssm", "action": "DescribeMaintenanceWindows"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "name", Require: plugin.Optional},
				{Name: "enabled", Require: plugin.Optional, Operators: []string{"=", "<>"}},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsSSMMaintenanceWindow,
				Tags: map[string]string{"service": "ssm", "action": "GetMaintenanceWindow"},
			},
			{
				Func: getAwsSSMMaintenanceWindowTags,
				Tags: map[string]string{"service": "ssm", "action": "ListTagsForResource"},
			},
			{
				Func: getMaintenanceWindowTargets,
				Tags: map[string]string{"service": "ssm", "action": "DescribeMaintenanceWindowTargets"},
			},
			{
				Func: getMaintenanceWindowTasks,
				Tags: map[string]string{"service": "ssm", "action": "DescribeMaintenanceWindowTasks"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ssmv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the Maintenance Window.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "window_id",
				Description: "The ID of the Maintenance Window.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enabled",
				Description: "Indicates whether the Maintenance Window is enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "allow_unassociated_targets",
				Description: "Indicates whether targets must be registered with the Maintenance Window before tasks can be defined for those targets.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsSSMMaintenanceWindow,
			},
			{
				Name:        "description",
				Description: "A description of the Maintenance Window.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the Maintenance Window",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSSMMaintenanceWindowTags,
				Transform:   transform.FromField("TagList"),
			},
			{
				Name:        "duration",
				Description: "The duration of the Maintenance Window in hours.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "cutoff",
				Description: "The number of hours before the end of the Maintenance Window that Systems Manager stops scheduling new tasks for execution.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "schedule",
				Description: "The schedule of the Maintenance Window in the form of a cron or rate expression.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "schedule_offset",
				Description: "The number of days to wait to run a Maintenance Window after the scheduled CRON expression date and time.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "targets",
				Description: "The targets of Maintenance Window.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMaintenanceWindowTargets,
				Transform:   transform.FromField("Targets"),
			},
			{
				Name:        "tasks",
				Description: "The Tasks of Maintenance Window.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMaintenanceWindowTasks,
				Transform:   transform.FromField("Tasks"),
			},
			{
				Name:        "created_date",
				Description: "The date the maintenance window was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getAwsSSMMaintenanceWindow,
			},
			{
				Name:        "end_date",
				Description: "The date and time, in ISO-8601 Extended format, for when the maintenance window is scheduled to become inactive. The maintenance window will not run after this specified time.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSSMMaintenanceWindow,
			},
			{
				Name:        "schedule_timezone",
				Description: "The schedule of the maintenance window in the form of a cron or rate expression.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSSMMaintenanceWindow,
			},
			{
				Name:        "start_date",
				Description: "The date and time, in ISO-8601 Extended format, for when the maintenance window is scheduled to become active.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "modified_date",
				Description: "The date the Maintenance Window was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getAwsSSMMaintenanceWindow,
			},
			// The value of NextExecutionTime is influenced by the timezone set by the user. Due to uncertainty regarding the date string's format, attempts to parse it into the time.RFC3339 format result in errors when the timezone isn't UTC. Consequently, we have designated the column type as a string.
			{
				Name:        "next_execution_time",
				Description: "The next time the maintenance window will actually run, taking into account any specified times for the Maintenance Window to become active or inactive.",
				Type:        proto.ColumnType_STRING,
			},

			/// Standard columns for all tables
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
				Hydrate:     getAwsSSMMaintenanceWindowTags,
				Transform:   transform.FromField("TagList").Transform(ssmMaintenanceWindowTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSSMMaintenanceWindowAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsSSMMaintenanceWindow(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := SSMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_maintenance_window.listAwsSSMMaintenanceWindow", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	maxItems := int32(100)
	input := &ssm.DescribeMaintenanceWindowsInput{}

	filters := buildSSMMaintenanceWindowFilter(d.Quals)

	if len(filters) > 0 {
		input.Filters = filters
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 10 {
				maxItems = int32(10)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	input.MaxResults = aws.Int32(maxItems)
	paginator := ssm.NewDescribeMaintenanceWindowsPaginator(svc, input, func(o *ssm.DescribeMaintenanceWindowsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ssm_maintenance_window.listAwsSSMMaintenanceWindow", "api_error", err)
			return nil, err
		}

		for _, parameter := range output.WindowIdentities {
			d.StreamListItem(ctx, parameter)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsSSMMaintenanceWindow(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var id string
	if h.Item != nil {
		id = *maintenanceWindowID(h.Item)
	} else {
		id = d.EqualsQuals["window_id"].GetStringValue()
	}

	// Empty id check
	if strings.TrimSpace(id) == "" {
		return nil, nil
	}

	// Create Session
	svc, err := SSMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_maintenance_window.getAwsSSMMaintenanceWindow", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &ssm.GetMaintenanceWindowInput{
		WindowId: &id,
	}

	// Get call
	data, err := svc.GetMaintenanceWindow(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_maintenance_window.getAwsSSMMaintenanceWindow", "api_error", err)
		return nil, err
	}

	return data, nil
}

func getMaintenanceWindowTargets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := maintenanceWindowID(h.Item)

	// Create Session
	svc, err := SSMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_maintenance_window.getMaintenanceWindowTargets", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &ssm.DescribeMaintenanceWindowTargetsInput{
		WindowId: id,
	}

	// Get call
	op, err := svc.DescribeMaintenanceWindowTargets(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_maintenance_window.getMaintenanceWindowTargets", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getMaintenanceWindowTasks(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := maintenanceWindowID(h.Item)

	// Create Session
	svc, err := SSMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_maintenance_window.getMaintenanceWindowTasks", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &ssm.DescribeMaintenanceWindowTasksInput{
		WindowId: id,
	}

	// Get call
	op, err := svc.DescribeMaintenanceWindowTasks(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_maintenance_window.getMaintenanceWindowTasks", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getAwsSSMMaintenanceWindowTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := maintenanceWindowID(h.Item)

	// Create Session
	svc, err := SSMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_maintenance_window.getAwsSSMMaintenanceWindowTags", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &ssm.ListTagsForResourceInput{
		ResourceType: types.ResourceTypeForTagging("MaintenanceWindow"),
		ResourceId:   id,
	}

	// Get call
	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_maintenance_window.getAwsSSMMaintenanceWindowTags", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getAwsSSMMaintenanceWindowAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	id := maintenanceWindowID(h.Item)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	aka := "arn:" + commonColumnData.Partition + ":ssm:" + region + ":" + commonColumnData.AccountId + ":maintenancewindow" + "/" + *id

	return []string{aka}, nil
}

/// TRANSFORM FUNCTIONS

func ssmMaintenanceWindowTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tagList := d.Value.([]types.Tag)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tagList != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
	} else {
		return nil, nil
	}

	return turbotTagsMap, nil
}

func maintenanceWindowID(item interface{}) *string {
	switch item := item.(type) {
	case *ssm.GetMaintenanceWindowOutput:
		return item.WindowId
	case types.MaintenanceWindowIdentity:
		return item.WindowId
	}
	return nil
}

//// UTILITY FUNCTION

// Build ssm maintenance window list call input filter
func buildSSMMaintenanceWindowFilter(quals plugin.KeyColumnQualMap) []types.MaintenanceWindowFilter {
	filters := make([]types.MaintenanceWindowFilter, 0)

	filterQuals := map[string]string{
		"name":    "Name",
		"enabled": "Enabled",
	}
	columnBool := []string{"enabled"}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := types.MaintenanceWindowFilter{
				Key: aws.String(filterName),
			}
			if slices.Contains(columnBool, columnName) {
				value := getQualsValueByColumn(quals, columnName, "boolean").(string)
				filter.Values = []string{cases.Title(language.English, cases.NoLower).String(value)}
			} else {
				value := getQualsValueByColumn(quals, columnName, "string")
				val, ok := value.(string)
				if ok {
					filter.Values = []string{val}
				} else {
					filter.Values = value.([]string)
				}
			}
			filters = append(filters, filter)
		}
	}
	return filters
}
