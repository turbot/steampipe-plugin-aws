package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsSSMMaintenanceWindow(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssm_maintenance_window",
		Description: "AWS SSM Maintenance Window",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("window_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"DoesNotExistException"}),
			},
			Hydrate: getAwsSSMMaintenanceWindow,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsSSMMaintenanceWindow,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "name", Require: plugin.Optional},
				{Name: "enabled", Require: plugin.Optional, Operators: []string{"=", "<>"}},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
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
			{
				Name:        "next_execution_time",
				Description: "The next time the maintenance window will actually run, taking into account any specified times for the Maintenance Window to become active or inactive.",
				Type:        proto.ColumnType_TIMESTAMP,
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
	plugin.Logger(ctx).Trace("listAwsSSMMaintenanceWindow")

	// Create session
	svc, err := SsmService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &ssm.DescribeMaintenanceWindowsInput{
		MaxResults: aws.Int64(100),
	}

	filters := buildSsmMentenanceWindowFilter(d.Quals)

	if len(filters) > 0 {
		input.Filters = filters
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 10 {
				input.MaxResults = aws.Int64(10)
			} else {
				input.MaxResults = limit
			}
		}
	}

	// List call
	err = svc.DescribeMaintenanceWindowsPages(
		input,
		func(page *ssm.DescribeMaintenanceWindowsOutput, isLast bool) bool {
			for _, parameter := range page.WindowIdentities {
				d.StreamListItem(ctx, parameter)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsSSMMaintenanceWindow(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsSSMMaintenanceWindow")

	var id string
	if h.Item != nil {
		id = *maintenanceWindowID(h.Item)
	} else {
		id = d.KeyColumnQuals["window_id"].GetStringValue()
	}

	// Create Session
	svc, err := SsmService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ssm.GetMaintenanceWindowInput{
		WindowId: &id,
	}

	// Get call
	data, err := svc.GetMaintenanceWindow(params)
	if err != nil {
		logger.Debug("getAwsSSMMaintenanceWindow", "ERROR", err)
		return nil, err
	}

	return data, nil
}

func getAwsSSMMaintenanceWindowAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsSSMMaintenanceWindowAkas")
	region := d.KeyColumnQualString(matrixKeyRegion)
	id := maintenanceWindowID(h.Item)
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	aka := "arn:" + commonColumnData.Partition + ":ssm:" + region + ":" + commonColumnData.AccountId + ":maintenancewindow" + "/" + *id

	return []string{aka}, nil
}

func getAwsSSMMaintenanceWindowTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsSSMMaintenanceWindowTags")

	id := maintenanceWindowID(h.Item)

	// Create Session
	svc, err := SsmService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ssm.ListTagsForResourceInput{
		ResourceType: types.String("MaintenanceWindow"),
		ResourceId:   id,
	}

	// Get call
	op, err := svc.ListTagsForResource(params)
	if err != nil {
		logger.Debug("getAwsSSMMaintenanceWindowTags", "ERROR", err)
		return nil, err
	}

	return op, nil
}

func getMaintenanceWindowTargets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getMaintenanceWindowTargets")

	id := maintenanceWindowID(h.Item)

	// Create Session
	svc, err := SsmService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ssm.DescribeMaintenanceWindowTargetsInput{
		WindowId: id,
	}

	// Get call
	op, err := svc.DescribeMaintenanceWindowTargets(params)
	if err != nil {
		logger.Debug("getMaintenanceWindowTargets", "ERROR", err)
		return nil, err
	}

	return op, nil
}

func getMaintenanceWindowTasks(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getMaintenanceWindowTasks")

	id := maintenanceWindowID(h.Item)

	// Create Session
	svc, err := SsmService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ssm.DescribeMaintenanceWindowTasksInput{
		WindowId: id,
	}

	// Get call
	op, err := svc.DescribeMaintenanceWindowTasks(params)
	if err != nil {
		logger.Debug("getMaintenanceWindowTasks", "ERROR", err)
		return nil, err
	}

	return op, nil
}

/// TRANSFORM FUNCTIONS

func ssmMaintenanceWindowTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("ssmMaintenanceWindowTagListToTurbotTags")
	tagList := d.Value.([]*ssm.Tag)

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
	case *ssm.MaintenanceWindowIdentity:
		return item.WindowId
	}
	return nil
}

//// UTILITY FUNCTION

// Build ssm maintenance window list call input filter
func buildSsmMentenanceWindowFilter(quals plugin.KeyColumnQualMap) []*ssm.MaintenanceWindowFilter {
	filters := make([]*ssm.MaintenanceWindowFilter, 0)

	filterQuals := map[string]string{
		"name":    "Name",
		"enabled": "Enabled",
	}
	columnBool := []string{"enabled"}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := ssm.MaintenanceWindowFilter{
				Key: aws.String(filterName),
			}
			if helpers.StringSliceContains(columnBool, columnName) {
				value := getQualsValueByColumn(quals, columnName, "boolean").(string)
				filter.Values = []*string{aws.String(cases.Title(language.English, cases.NoLower).String(value))}
			} else {
				value := getQualsValueByColumn(quals, columnName, "string")
				val, ok := value.(string)
				if ok {
					filter.Values = []*string{&val}
				} else {
					filter.Values = value.([]*string)
				}
			}
			filters = append(filters, &filter)
		}
	}
	return filters
}
