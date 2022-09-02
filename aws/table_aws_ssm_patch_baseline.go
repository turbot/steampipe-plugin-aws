package aws

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
)

//// TABLE DEFINITION

func tableAwsSSMPatchBaseline(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssm_patch_baseline",
		Description: "AWS SSM Patch Baseline",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("baseline_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"DoesNotExistException", "InvalidResourceId", "InvalidParameter", "ValidationException"}),
			},
			Hydrate: getPatchBaseline,
		},
		List: &plugin.ListConfig{
			Hydrate: describePatchBaselines,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "name", Require: plugin.Optional},
				{Name: "operating_system", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the patch baseline.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "baseline_id",
				Description: "The ID of the retrieved patch baseline.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel().Transform(lastPathElement),
			},
			{
				Name:        "description",
				Description: "A description of the patch baseline.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "operating_system",
				Description: "Returns the operating system specified for the patch baseline.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_date",
				Description: "The date the patch baseline was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getPatchBaseline,
			},
			{
				Name:        "modified_date",
				Description: "The date the patch baseline was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getPatchBaseline,
			},
			{
				Name:        "approved_patches_compliance_level",
				Description: "Returns the specified compliance severity level for approved patches in the patch baseline.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getPatchBaseline,
			},
			{
				Name:        "approved_patches_enable_non_security",
				Description: "Indicates whether the list of approved patches includes non-security updates that should be applied to the instances. The default value is 'false'. Applies to Linux instances only.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getPatchBaseline,
			},
			{
				Name:        "approval_rules",
				Description: "A set of rules used to include patches in the baseline.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getPatchBaseline,
			},
			{
				Name:        "approved_patches",
				Description: "A list of explicitly approved patches for the baseline.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getPatchBaseline,
			},
			{
				Name:        "global_filters",
				Description: "A set of global filters used to exclude patches from the baseline.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getPatchBaseline,
			},
			{
				Name:        "patch_groups",
				Description: "Patch groups included in the patch baseline.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getPatchBaseline,
			},
			{
				Name:        "rejected_patches_action",
				Description: "The action specified to take on patches included in the RejectedPatches list. A patch can be allowed only if it is a dependency of another package, or blocked entirely along with packages that include it as a dependency.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getPatchBaseline,
			},
			{
				Name:        "rejected_patches",
				Description: "A list of explicitly rejected patches for the baseline.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getPatchBaseline,
			},
			{
				Name:        "sources",
				Description: "Information about the patches to use to update the instances, including target operating systems and source repositories. Applies to Linux instances only.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getPatchBaseline,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the patch baseline.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSSMPatchBaselineTags,
				Transform:   transform.FromField("TagList"),
			},

			// Standard columns for all tables
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
				Hydrate:     getAwsSSMPatchBaselineTags,
				Transform:   transform.FromField("TagList").Transform(ssmTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsSSMPatchBaselineAkas,
				Transform:   transform.FromValue(),
			},
			// TODO: The below mention field is coming from list call, but not from get call.
			// Need to check, if there is another way to fetch this value.

			// {
			// 	Name:        "default_baseline",
			// 	Description: "Whether this is the default baseline. Note that Systems Manager supports creating multiple default patch baselines. For example, you can create a default patch baseline for each operating system.",
			// 	Type:        proto.ColumnType_BOOL,
			// },
		}),
	}
}

//// LIST FUNCTION

func describePatchBaselines(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("describePatchBaselines")

	// Create session
	svc, err := SsmService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	// Adding a filter to filter out all the predefined patch baseline, that does not belongs to the user or account
	params := &ssm.DescribePatchBaselinesInput{
		MaxResults: aws.Int64(100),
	}

	ownerFilter := &ssm.PatchOrchestratorFilter{
		Key:    aws.String("OWNER"),
		Values: []*string{aws.String("Self")},
	}

	filters := append(buildSsmPatchBaselineFilter(d.Quals), ownerFilter)
	params.Filters = filters

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *params.MaxResults {
			if *limit < 1 {
				params.MaxResults = aws.Int64(1)
			} else {
				params.MaxResults = limit
			}
		}
	}

	// List call
	err = svc.DescribePatchBaselinesPages(
		params,
		func(page *ssm.DescribePatchBaselinesOutput, isLast bool) bool {
			for _, baseline := range page.BaselineIdentities {
				var rowData *ssm.GetPatchBaselineOutput
				if baseline != nil {
					rowData = &ssm.GetPatchBaselineOutput{
						BaselineId:      baseline.BaselineId,
						Name:            baseline.BaselineName,
						OperatingSystem: baseline.OperatingSystem,
						Description:     baseline.BaselineDescription,
					}
				}
				d.StreamListItem(ctx, rowData)

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

func getPatchBaseline(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getPatchBaseline")

	var baselineID string
	if h.Item != nil {
		baselineID = *h.Item.(*ssm.GetPatchBaselineOutput).BaselineId
	} else {
		quals := d.KeyColumnQuals
		baselineID = quals["baseline_id"].GetStringValue()
	}

	// get service
	svc, err := SsmService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ssm.GetPatchBaselineInput{
		BaselineId: aws.String(baselineID),
	}

	// Get call
	data, err := svc.GetPatchBaseline(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getPatchBaseline__", "ERROR", err)
		return nil, err
	}
	return data, nil
}

// API call for fetching tag list
func getAwsSSMPatchBaselineTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsSSMPatchBaselineTags")

	baseline := h.Item.(*ssm.GetPatchBaselineOutput)

	// Create Session
	svc, err := SsmService(ctx, d)
	if err != nil {
		return nil, err
	}

	baselineIDSplitted := strings.Split(*baseline.BaselineId, "/")
	id := baselineIDSplitted[len(baselineIDSplitted)-1]

	// Build the params
	params := &ssm.ListTagsForResourceInput{
		ResourceType: types.String("PatchBaseline"),
		ResourceId:   &id,
	}

	// Get call
	op, err := svc.ListTagsForResource(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getAwsSSMPatchBaselineTags", "ERROR", err)
		return nil, err
	}

	return op, nil
}

func getAwsSSMPatchBaselineAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsSSMPatchBaselineAkas")
	parameterData := h.Item.(*ssm.GetPatchBaselineOutput)
	region := d.KeyColumnQualString(matrixKeyRegion)

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	c, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)

	aka := "arn:" + commonColumnData.Partition + ":ssm:" + region + ":" + commonColumnData.AccountId + ":patchbaseline"

	if strings.HasPrefix(*parameterData.BaselineId, "/") {
		aka = aka + *parameterData.BaselineId
	} else {
		aka = aka + "/" + *parameterData.BaselineId
	}

	return []string{aka}, nil
}

//// UTILITY FUNCTION

// Build ssm patch baseline list call input filter
func buildSsmPatchBaselineFilter(quals plugin.KeyColumnQualMap) []*ssm.PatchOrchestratorFilter {
	filters := make([]*ssm.PatchOrchestratorFilter, 0)

	filterQuals := map[string]string{
		"name":             "NAME_PREFIX",
		"operating_system": "OPERATING_SYSTEM",
	}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := ssm.PatchOrchestratorFilter{
				Key: aws.String(filterName),
			}

			value := getQualsValueByColumn(quals, columnName, "string")
			val, ok := value.(string)
			if ok {
				filter.Values = []*string{&val}
			} else {
				filter.Values = value.([]*string)
			}
			filters = append(filters, &filter)
		}
	}
	return filters
}
