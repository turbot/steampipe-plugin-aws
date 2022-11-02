package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"

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
	// Create session
	svc, err := ResourceGroupsTaggingClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_tagging_resource.listTaggingResources", "get_client_error", err)
		return nil, err
	}

	input := &resourcegroupstaggingapi.GetResourcesInput{
		ResourcesPerPage: aws.Int32(100),
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
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_tagging_resource.listTaggingResources", "api_error", err)
			return nil, err
		}

		for _, resource := range output.ResourceTagMappingList {
			d.StreamListItem(ctx, resource)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getTaggingResource(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	arn := d.KeyColumnQuals["arn"].GetStringValue()

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
