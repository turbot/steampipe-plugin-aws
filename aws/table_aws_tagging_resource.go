package aws

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
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
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:       "resource_types",
					Require:    plugin.Optional,
					Operators:  []string{"="},
					CacheMatch: query_cache.CacheMatchExact,
				},
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
				Name:        "resource_types",
				Description: "The resource types to filter by in the form of an array of strings.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromQual("resource_types"),
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

	// Parse resource type filters
	var allResourceTypes []string
	resource_types := d.EqualsQuals["resource_types"].GetJsonbValue()
	if resource_types != "" {
		err := json.Unmarshal([]byte(resource_types), &allResourceTypes)
		if err != nil {
			return nil, errors.New("failed to parse 'resource_types' qualifier: value must be a JSON array of strings, e.g. [\"ec2:instance\", \"s3:bucket\", \"rds\"]")
		}
	}

	// Split resource types into batches that don't exceed 100 items per request
	resourceTypeBatches := batchResourceTypes(allResourceTypes, 100)

	// If no resource types specified, make a single request
	if len(resourceTypeBatches) == 0 {
		resourceTypeBatches = append(resourceTypeBatches, []string{})
	}

	// Track seen resources to avoid duplicates across batches
	seenResources := make(map[string]bool)

	// Process each batch
	for _, resourceTypesBatch := range resourceTypeBatches {
		err := fetchResourcesForBatch(ctx, d, svc, resourceTypesBatch, seenResources)
		if err != nil {
			return nil, err
		}

		// Check if we've hit the limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

// fetchResourcesForBatch fetches resources for a specific batch of resource types
func fetchResourcesForBatch(ctx context.Context, d *plugin.QueryData, svc *resourcegroupstaggingapi.Client, resourceTypes []string, seenResources map[string]bool) error {
	input := &resourcegroupstaggingapi.GetResourcesInput{
		ResourcesPerPage: aws.Int32(100),
	}

	if len(resourceTypes) > 0 {
		input.ResourceTypeFilters = resourceTypes
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
			return err
		}

		for _, resource := range output.ResourceTagMappingList {
			// Deduplicate based on ARN
			arn := aws.ToString(resource.ResourceARN)
			if seenResources[arn] {
				continue // Skip duplicate
			}
			seenResources[arn] = true

			d.StreamListItem(ctx, resource)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil
			}
		}
	}

	return nil
}

// batchResourceTypes splits resource types into batches that don't exceed maxItems per request
func batchResourceTypes(resourceTypes []string, maxItems int) [][]string {
	if len(resourceTypes) == 0 {
		return [][]string{}
	}

	var batches [][]string
	var currentBatch []string
	currentItems := 0

	for _, resourceType := range resourceTypes {
		// If adding this would exceed the limit, start a new batch
		if currentItems+1 > maxItems {
			batches = append(batches, currentBatch)
			currentBatch = []string{resourceType}
			currentItems = 1
		} else {
			currentBatch = append(currentBatch, resourceType)
			currentItems++
		}
	}

	// Add the last batch if it has items
	if len(currentBatch) > 0 {
		batches = append(batches, currentBatch)
	}

	return batches
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
