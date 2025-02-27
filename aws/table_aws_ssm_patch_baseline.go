package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSSMPatchBaseline(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssm_patch_baseline",
		Description: "AWS SSM Patch Baseline",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("baseline_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"DoesNotExistException", "InvalidResourceId", "InvalidParameter", "ValidationException"}),
			},
			Hydrate: getPatchBaseline,
			Tags:    map[string]string{"service": "ssm", "action": "GetPatchBaseline"},
		},
		List: &plugin.ListConfig{
			Hydrate: describePatchBaselines,
			Tags:    map[string]string{"service": "ssm", "action": "DescribePatchBaselines"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "name", Require: plugin.Optional},
				{Name: "operating_system", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getPatchBaseline,
				Tags: map[string]string{"service": "ssm", "action": "GetPatchBaseline"},
			},
			{
				Func: getAwsSSMPatchBaselineTags,
				Tags: map[string]string{"service": "ssm", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_SSM_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the patch baseline.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BaselineName", "Name"),
			},
			{
				Name:        "baseline_id",
				Description: "The ID of the retrieved patch baseline.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BaselineId").Transform(lastPathElement),
			},
			{
				Name:        "description",
				Description: "A description of the patch baseline.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description", "BaselineDescription"),
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
				Transform:   transform.FromField("BaselineName", "Name"),
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

	// Create session
	svc, err := SSMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_patch_baseline.describePatchBaselines", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	// Adding a filter to filter out all the predefined patch baseline, that does not belongs to the user or account
	maxItems := int32(100)
	input := &ssm.DescribePatchBaselinesInput{}

	ownerFilter := types.PatchOrchestratorFilter{
		Key:    aws.String("OWNER"),
		Values: []string{"Self"},
	}

	filters := append(buildSSMPatchBaselineFilter(d.Quals), ownerFilter)
	input.Filters = filters

	// Reduce the basic request limit down if the user has only requested a small number of rows
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
	paginator := ssm.NewDescribePatchBaselinesPaginator(svc, input, func(o *ssm.DescribePatchBaselinesPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ssm_patch_baseline.describePatchBaselines", "api_error", err)
			return nil, err
		}

		for _, baseline := range output.BaselineIdentities {
			d.StreamListItem(ctx, baseline)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getPatchBaseline(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var baselineID string
	if h.Item != nil {
		baselineID = *h.Item.(types.PatchBaselineIdentity).BaselineId
	} else {
		quals := d.EqualsQuals
		baselineID = quals["baseline_id"].GetStringValue()
	}

	// Empty baseline id check
	if strings.TrimSpace(baselineID) == "" {
		return nil, nil
	}

	// get service
	svc, err := SSMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_patch_baseline.getPatchBaseline", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build the params
	params := &ssm.GetPatchBaselineInput{
		BaselineId: aws.String(baselineID),
	}

	// Get call
	data, err := svc.GetPatchBaseline(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_patch_baseline.getPatchBaseline", "api_error", err)
		return nil, err
	}
	return data, nil
}

// API call for fetching tag list
func getAwsSSMPatchBaselineTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	baselineId := getPatchBaselineID(h.Item)

	// Create Session
	svc, err := SSMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_patch_baseline.getAwsSSMPatchBaselineTags", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	baselineIDSplitted := strings.Split(baselineId, "/")
	id := baselineIDSplitted[len(baselineIDSplitted)-1]

	// Build the params
	params := &ssm.ListTagsForResourceInput{
		ResourceType: types.ResourceTypeForTagging("PatchBaseline"),
		ResourceId:   &id,
	}

	// Get call
	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_patch_baseline.getAwsSSMPatchBaselineTags", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getAwsSSMPatchBaselineAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	baselineId := getPatchBaselineID(h.Item)
	region := d.EqualsQualString(matrixKeyRegion)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_patch_baseline.getAwsSSMPatchBaselineAkas", "common_data_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)

	aka := "arn:" + commonColumnData.Partition + ":ssm:" + region + ":" + commonColumnData.AccountId + ":patchbaseline"

	if strings.HasPrefix(baselineId, "/") {
		aka = aka + baselineId
	} else {
		aka = aka + "/" + baselineId
	}

	return []string{aka}, nil
}

//// UTILITY FUNCTION

// Build ssm patch baseline list call input filter
func buildSSMPatchBaselineFilter(quals plugin.KeyColumnQualMap) []types.PatchOrchestratorFilter {
	filters := make([]types.PatchOrchestratorFilter, 0)

	filterQuals := map[string]string{
		"name":             "NAME_PREFIX",
		"operating_system": "OPERATING_SYSTEM",
	}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := types.PatchOrchestratorFilter{
				Key: aws.String(filterName),
			}

			value := getQualsValueByColumn(quals, columnName, "string")
			val, ok := value.(string)
			if ok {
				filter.Values = []string{val}
			} else {
				filter.Values = value.([]string)
			}
			filters = append(filters, filter)
		}
	}
	return filters
}

func getPatchBaselineID(data any) string {
	switch item := data.(type) {
	case *ssm.GetPatchBaselineOutput:
		return *item.BaselineId
	case types.PatchBaselineIdentity:
		return *item.BaselineId
	}
	return ""
}
