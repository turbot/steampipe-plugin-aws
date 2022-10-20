package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go-v2/service/redshiftserverless"
	"github.com/aws/aws-sdk-go-v2/service/redshiftserverless/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsRedshiftServerlessWorkgroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_redshiftserverless_workgroup",
		Description: "AWS Redshift Serverless Workgroup",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("workgroup_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getRedshiftServerlessWorkgroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listRedshiftServerlessWorkgroups,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "workgroup_name",
				Description: "The name of the workgroup.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "workgroup_id",
				Description: "The unique identifier of the workgroup.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "workgroup_arn",
				Description: "The Amazon Resource Name (ARN) that links to the workgroup.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the workgroup.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "base_capacity",
				Description: "The base data warehouse capacity of the workgroup in Redshift Processing Units (RPUs).",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "creation_date",
				Description: "The creation date of the workgroup.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "enhanced_vpc_routing",
				Description: "The value that specifies whether to enable enhanced virtual private cloud (VPC) routing, which forces Amazon Redshift Serverless to route traffic through your VPC.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "namespace_name",
				Description: "The namespace the workgroup is associated with.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "publicly_accessible",
				Description: "A value that specifies whether the workgroup can be accessible from a public network.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "config_parameters",
				Description: "An array of parameters to set for finer control over a database.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "endpoint",
				Description: "The endpoint that is created from the workgroup.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "security_group_ids",
				Description: "An array of security group IDs to associate with the workgroup.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "subnet_ids",
				Description: "An array of subnet IDs the workgroup is associated with.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "The list of tags for the workgroup.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWorkgroupTags,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkgroupName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWorkgroupTags,
				Transform:   transform.From(getWorkgroupTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WorkgroupArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listRedshiftServerlessWorkgroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create Session
	svc, err := RedshiftServerlessClient(ctx, d)
	if err != nil {
		logger.Error("aws_redshiftserverless_workgroup.listRedshiftServerlessWorkgroups", "service_creation_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	input := &redshiftserverless.ListWorkgroupsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := redshiftserverless.NewListWorkgroupsPaginator(svc, input, func(o *redshiftserverless.ListWorkgroupsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_redshiftserverless_workgroup.listRedshiftServerlessWorkgroups", "api_error", err)
			return nil, err
		}

		for _, workgroup := range output.Workgroups {
			d.StreamListItem(ctx, workgroup)
		}

		// Context may get cancelled due to manual cancellation or if the limit has been reached
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getRedshiftServerlessWorkgroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	name := d.KeyColumnQuals["workgroup_name"].GetStringValue()
	// Return nil, if no input provided
	if name == "" {
		return nil, nil
	}

	// Create service
	svc, err := RedshiftServerlessClient(ctx, d)
	if err != nil {
		logger.Error("aws_redshiftserverless_workgroup.getRedshiftServerlessWorkgroup", "service_creation_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build params
	params := &redshiftserverless.GetWorkgroupInput{
		WorkgroupName: aws.String(name),
	}

	op, err := svc.GetWorkgroup(ctx, params)
	if err != nil {
		logger.Error("aws_redshiftserverless_workgroup.getRedshiftServerlessWorkgroup", "api_error", err)
		return nil, err
	}
	return *op.Workgroup, nil
}

func getWorkgroupTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	arn := *h.Item.(types.Workgroup).WorkgroupArn

	// Create service
	svc, err := RedshiftServerlessClient(ctx, d)
	if err != nil {
		logger.Error("aws_redshiftserverless_workgroup.getWorkgroupTags", "service_creation_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build params
	params := &redshiftserverless.ListTagsForResourceInput{
		ResourceArn: aws.String(arn),
	}

	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		logger.Error("aws_redshiftserverless_workgroup.getWorkgroupTags", "api_error", err)
		return nil, err
	}
	return op, nil
}

//// TRANSFORM FUNCTIONS

func getWorkgroupTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	op := d.HydrateItem.(*redshiftserverless.ListTagsForResourceOutput)

	if op.Tags != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range op.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
