package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lakeformation"

	lakeformationv1 "github.com/aws/aws-sdk-go/service/lakeformation"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsLakeformationResource(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_lakeformation_resource",
		Description: "AWS Lake Formation Resource",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("resource_arn"),
			Hydrate:    getLakeformationResource,
			Tags:       map[string]string{"service": "lakeformation", "action": "DescribeResource"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"EntityNotFoundException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listLakeformationResources,
			Tags:    map[string]string{"service": "lakeformation", "action": "ListResources"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(lakeformationv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "resource_arn",
				Description: "The Amazon Resource Name (ARN) of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "role_arn",
				Description: "The IAM role that registered a resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "hybrid_access_enabled",
				Description: "Indicates whether the data access of tables pointing to the location can be managed by both Lake Formation permissions as well as Amazon S3 bucket policies.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "last_modified",
				Description: "The date and time the resource was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "with_federation",
				Description: "Whether or not the resource is a federated resource.",
				Type:        proto.ColumnType_BOOL,
			},

			/// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ResourceArn"),
			},
		}),
	}
}

//// LIST FUNCTION

func listLakeformationResources(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := LakeFormationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lakeformation_resource.listLakeformationResources", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	maxItems := int32(1000)
	input := &lakeformation.ListResourcesInput{}

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
	paginator := lakeformation.NewListResourcesPaginator(svc, input, func(o *lakeformation.ListResourcesPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_lakeformation_resource.listLakeformationResources", "api_error", err)
			return nil, err
		}

		for _, resource := range output.ResourceInfoList {
			d.StreamListItem(ctx, resource)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getLakeformationResource(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	arn := d.EqualsQuals["resource_arn"].GetStringValue()

	// Empty id check
	if arn == "" {
		return nil, nil
	}

	// Create session
	svc, err := LakeFormationClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lakeformation_resource.getLakeformationResource", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &lakeformation.DescribeResourceInput{
		ResourceArn: &arn,
	}

	res, err := svc.DescribeResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_lakeformation_resource.getLakeformationResource", "api_error", err)
		return nil, err
	}

	if res != nil && res.ResourceInfo != nil {
		// The DescribeResource API does not return the ResourceArn property, even though it is present in the response attributes.
		// The AWS CLI exhibits the same behavior as the API.
		// Therefore, we are assigning the value from the qualifiers (quals) here.
		res.ResourceInfo.ResourceArn = &arn

		return *res.ResourceInfo, nil
	}

	return nil, nil
}
