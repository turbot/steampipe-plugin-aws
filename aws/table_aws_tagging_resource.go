package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsTaggingResource(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_tagging_resource",
		Description: "AWS Tagging Resource",
		Get: &plugin.GetConfig{
			Hydrate:           getTaggingResource,
			KeyColumns:        plugin.SingleColumn("arn"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidParameterException"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listTaggingResources,
		},
		GetMatrixItem: BuildRegionList,
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
	svc, err := TaggignResourceService(ctx, d)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.GetResourcesPages(
		&resourcegroupstaggingapi.GetResourcesInput{},
		func(page *resourcegroupstaggingapi.GetResourcesOutput, isLast bool) bool {
			for _, resource := range page.ResourceTagMappingList {
				d.StreamListItem(ctx, resource)
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
	svc, err := TaggignResourceService(ctx, d)
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
