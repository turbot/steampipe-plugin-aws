package aws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v6/query_cache"
)

func tableAwsTaggingResource(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_tagging_resource",
		Description: "AWS Tagging Resource",
		Get: &plugin.GetConfig{
			Hydrate:    getTaggingResource,
			Tags:       map[string]string{"service": "tag", "action": "GetResources"},
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listTaggingResources,
			Tags:    map[string]string{"service": "tag", "action": "GetResources"},
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "tag_filter", Operators: []string{"="}, Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_TAGGING_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceARN").Transform(arnToTitle),
			},
			{
				Name:        "arn",
				Description: "The ARN of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceARN"),
			},
			{
				Name:        "compliance_status",
				Description: "Whether a resource is compliant with the effective tag policy.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ComplianceDetails.ComplianceStatus"),
			},
			{
				Name:        "keys_with_noncompliant_values",
				Description: "These are keys defined in the effective policy that are on the resource with either incorrect case treatment or noncompliant values.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ComplianceDetails.KeysWithNoncompliantValues"),
			},
			{
				Name:        "noncompliant_keys",
				Description: "These tag keys on the resource are noncompliant with the effective tag policy.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ComplianceDetails.NoncompliantKeys"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the parameter.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "tag_filter",
				Description: "A list of TagFilters used to filter resources by tags. Specify a JSON array of objects with 'key' and optional 'values' fields, e.g., [{\"key\":\"Environment\",\"values\":[\"prod\",\"dev\"]}].",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromQual("tag_filter"),
			},

			/// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceARN").Transform(arnToTitle),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags").Transform(resourceTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ResourceARN").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listTaggingResources(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := ResourceGroupsTaggingClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_tagging_resource.listTaggingResources", "get_client_error", err)
		return nil, err
	}

	input := &resourcegroupstaggingapi.GetResourcesInput{
		ResourcesPerPage: aws.Int32(100),
	}

	// Build tag filters from quals
	tagFilters, err := buildTagFilter(d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_tagging_resource.listTaggingResources", "build_tag_filter_error", err)
		return nil, err
	}
	if len(tagFilters) > 0 {
		input.TagFilters = tagFilters
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < *input.ResourcesPerPage {
			if limit < 1 {
				input.ResourcesPerPage = aws.Int32(1)
			} else {
				input.ResourcesPerPage = aws.Int32(limit)
			}
		}
	}

	paginator := resourcegroupstaggingapi.NewGetResourcesPaginator(svc, input, func(o *resourcegroupstaggingapi.GetResourcesPaginatorOptions) {
		o.Limit = *input.ResourcesPerPage
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_tagging_resource.listTaggingResources", "api_error", err)
			return nil, err
		}

		for _, resource := range output.ResourceTagMappingList {
			d.StreamListItem(ctx, resource)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getTaggingResource(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	arn := d.EqualsQuals["arn"].GetStringValue()

	// Create session
	svc, err := ResourceGroupsTaggingClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_tagging_resource.getTaggingResource", "get_client_error", err)
		return nil, err
	}

	param := &resourcegroupstaggingapi.GetResourcesInput{
		ResourceARNList: []string{arn},
	}

	op, err := svc.GetResources(ctx, param)
	if err != nil {
		plugin.Logger(ctx).Error("aws_tagging_resource.getTaggingResource", "api_error", err)
		return nil, err
	}

	if op != nil && len(op.ResourceTagMappingList) > 0 {
		return op.ResourceTagMappingList[0], nil
	}

	return nil, nil
}

//// HELPER FUNCTIONS

// buildTagFilter constructs TagFilter objects from the tag_filter qual
func buildTagFilter(d *plugin.QueryData) ([]types.TagFilter, error) {
	var tagFilters []types.TagFilter

	if d.Quals["tag_filter"] != nil {
		for _, q := range d.Quals["tag_filter"].Quals {
			val := q.Value.GetJsonbValue()
			if val != "" && q.Operator == "=" {
				// Parse JSON: [{"key":"environment","values":["prod","dev"]},{"key":"team"}]
				var filterInput []map[string]interface{}
				if err := json.Unmarshal([]byte(val), &filterInput); err != nil {
					return nil, err
				}

				for _, filter := range filterInput {
					keyStr, keyOk := filter["key"].(string)
					if !keyOk || keyStr == "" {
						continue
					}

					tagFilter := types.TagFilter{
						Key: aws.String(keyStr),
					}

				// Values are optional
				if valuesRaw, exists := filter["values"]; exists {
					valuesArray, ok := valuesRaw.([]interface{})
					if !ok {
						return nil, fmt.Errorf("tag_filter: 'values' must be an array of strings, got %T", valuesRaw)
					}
					var values []string
					for _, v := range valuesArray {
						if vStr, vOk := v.(string); vOk && vStr != "" {
							values = append(values, vStr)
						}
					}
					if len(values) > 0 {
						tagFilter.Values = values
					}
				}

					tagFilters = append(tagFilters, tagFilter)
				}
			}
		}
	}

	return tagFilters, nil
}

//// TRANSFORM FUNCTIONS

func resourceTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tagList := d.Value.([]types.Tag)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tagList != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}
