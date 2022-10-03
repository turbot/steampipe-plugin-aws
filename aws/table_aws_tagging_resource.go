package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsTaggingResource(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_tagging_resource",
		Description: "AWS Tagging Resource",
		Get: &plugin.GetConfig{
			Hydrate:    getTaggingResource,
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"InvalidParameterException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listTaggingResources,
		},
		GetMatrixItemFunc: BuildRegionList,
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
	plugin.Logger(ctx).Trace("listTaggingResources")

	// Create session
	svc, err := TaggingResourceService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &resourcegroupstaggingapi.GetResourcesInput{
		ResourcesPerPage: aws.Int64(100),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.ResourcesPerPage {
			if *limit < 1 {
				input.ResourcesPerPage = aws.Int64(1)
			} else {
				input.ResourcesPerPage = limit
			}
		}
	}

	// List call
	err = svc.GetResourcesPages(
		input,
		func(page *resourcegroupstaggingapi.GetResourcesOutput, isLast bool) bool {
			for _, resource := range page.ResourceTagMappingList {
				d.StreamListItem(ctx, resource)

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

func getTaggingResource(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getTaggingResource")

	arn := d.KeyColumnQuals["arn"].GetStringValue()

	// Create session
	svc, err := TaggingResourceService(ctx, d)
	if err != nil {
		return nil, err
	}

	param := &resourcegroupstaggingapi.GetResourcesInput{
		ResourceARNList: []*string{&arn},
	}

	op, err := svc.GetResources(param)
	if err != nil {
		return nil, err
	}

	if op != nil && len(op.ResourceTagMappingList) > 0 {
		return op.ResourceTagMappingList[0], nil
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS

func resourceTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("resourceTagListToTurbotTags")
	tagList := d.Value.([]*resourcegroupstaggingapi.Tag)

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
