package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

func tableAwsCloudControlResource(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudcontrol_resource",
		Description: "AWS Cloud Control Resource",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "type_name", Require: plugin.Required},
				{Name: "resource_model", Require: plugin.Optional},
			},
			Hydrate: listCloudControlResources,
			Tags:    map[string]string{"service": "cloudformation", "action": "ListResources"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "type_name", CacheMatch: query_cache.CacheMatchExact},
				{Name: "identifier"},
			},
			Hydrate: getCloudControlResource,
			Tags:    map[string]string{"service": "cloudformation", "action": "GetResource"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getCloudControlResource,
				Tags: map[string]string{"service": "cloudformation", "action": "GetResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_CLOUDCONTROLAPI_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "type_name",
				Description: "The name of the resource type.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("type_name"),
			},
			{
				Name:        "identifier",
				Description: "The identifier for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_model",
				Description: "The resource model to use to select the resources to return.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("resource_model"),
			},
			{
				Name:        "properties",
				Description: "Represents information about a provisioned resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudControlResource,
			},
		}),
	}
}

type cloudControlResource struct {
	Identifier *string
	Properties *string
}

//// LIST FUNCTION

func listCloudControlResources(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := CloudControlClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudcontrol_resource.listCloudControlResources", "connection_error", err)
		return nil, err
	}

	typeName := d.EqualsQuals["type_name"].GetStringValue()
	resourceModel := d.EqualsQuals["resource_model"].GetStringValue()
	if strings.TrimSpace(typeName) == "" {
		return nil, nil
	}

	// Set MaxResults to the maximum number allowed
	maxItems := int32(100)
	input := &cloudcontrol.ListResourcesInput{
		TypeName: aws.String(typeName),
	}

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

	if len(resourceModel) > 0 {
		input.ResourceModel = aws.String(resourceModel)
	}

	paginator := cloudcontrol.NewListResourcesPaginator(svc, input, func(o *cloudcontrol.ListResourcesPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudcontrol_resource.listCloudControlResources", "api_error", err)
			return nil, err
		}

		for _, resource := range output.ResourceDescriptions {
			d.StreamListItem(ctx, resource)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCloudControlResource(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := CloudControlClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudcontrol_resource.getCloudControlResource", "connection_error", err)
		return nil, err
	}

	var identifier string
	typeName := d.EqualsQuals["type_name"].GetStringValue()

	if h.Item != nil {
		resource := h.Item.(types.ResourceDescription)
		identifier = *resource.Identifier
		resourceProperties := *resource.Properties

		// S3 buckets are too expensive to hydrate, so just return the list
		// properties
		if typeName == "AWS::S3::Bucket" {
			return types.ResourceDescription{
				Identifier: aws.String(identifier),
				Properties: aws.String(resourceProperties),
			}, nil
		}
	} else {
		identifier = d.EqualsQuals["identifier"].GetStringValue()
	}

	input := &cloudcontrol.GetResourceInput{
		Identifier: aws.String(identifier),
		TypeName:   aws.String(typeName),
	}

	item, err := svc.GetResource(ctx, input)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudcontrol_resource.getCloudControlResource", "api_error", err)
		return nil, err
	}

	properties := item.ResourceDescription.Properties

	return &cloudControlResource{
		Identifier: aws.String(identifier),
		Properties: properties,
	}, nil
}
