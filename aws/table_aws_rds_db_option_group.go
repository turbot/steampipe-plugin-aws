package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"

	rdsv1 "github.com/aws/aws-sdk-go/service/rds"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsRDSDBOptionGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_rds_db_option_group",
		Description: "AWS RDS DB Option Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"OptionGroupNotFoundFault"}),
			},
			Hydrate: getRDSDBOptionGroup,
			Tags:    map[string]string{"service": "rds", "action": "DescribeOptionGroups"},
		},
		List: &plugin.ListConfig{
			Hydrate: listRDSDBOptionGroups,
			Tags:    map[string]string{"service": "rds", "action": "DescribeOptionGroups"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "engine_name", Require: plugin.Optional},
				{Name: "major_engine_version", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsRDSOptionGroupTags,
				Tags: map[string]string{"service": "rds", "action": "ListTagsForResource"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(rdsv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name to identify the option group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OptionGroupName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the option group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OptionGroupArn"),
			},
			{
				Name:        "description",
				Description: "Provides a description of the option group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OptionGroupDescription"),
			},
			{
				Name:        "allows_vpc_and_non_vpc_instance_memberships",
				Description: "Specifies whether this option group can be applied to both VPC and non-VPC instances.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "engine_name",
				Description: "Indicates the name of the engine that this option group can be applied to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "major_engine_version",
				Description: "Indicates the major engine version associated with this option group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "Indicates the ID of the VPC, option group can be applied.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "options",
				Description: "Indicates what options are available in the option group.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached to the option group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsRDSOptionGroupTags,
				Transform:   transform.FromField("TagList"),
			},

			// Standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsRDSOptionGroupTags,
				Transform:   transform.From(getRDSDBOptionGroupTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OptionGroupName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("OptionGroupArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listRDSDBOptionGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	// Create Session
	svc, err := RDSClient(ctx, d)
	if err != nil {
		logger.Error("aws_rds_db_option_group.listRDSDBOptionGroups", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 20 {
				maxLimit = 20
			} else {
				maxLimit = limit
			}
		}
	}

	input := &rds.DescribeOptionGroupsInput{
		MaxRecords: &maxLimit,
	}

	equalQuals := d.EqualsQuals
	if equalQuals["engine_name"] != nil {
		input.EngineName = aws.String(equalQuals["engine_name"].GetStringValue())
	}

	// We must need to pass engine name if we are passing the major engine version
	if equalQuals["engine_name"] != nil && equalQuals["major_engine_version"] != nil {
		input.MajorEngineVersion = aws.String(equalQuals["major_engine_version"].GetStringValue())
	}

	paginator := rds.NewDescribeOptionGroupsPaginator(svc, input, func(o *rds.DescribeOptionGroupsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			logger.Error("aws_rds_db_option_group.listRDSDBOptionGroups", "api_error", err)
			return nil, err
		}

		for _, items := range output.OptionGroupsList {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getRDSDBOptionGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()

	// Create service
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_option_group.getRDSDBOptionGroup", "connection_error", err)
		return nil, err
	}

	params := &rds.DescribeOptionGroupsInput{
		OptionGroupName: aws.String(name),
	}

	op, err := svc.DescribeOptionGroups(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_option_group.getRDSDBOptionGroup", "api_error", err)
		return nil, err
	}

	if op.OptionGroupsList != nil && len(op.OptionGroupsList) > 0 {
		return op.OptionGroupsList[0], nil
	}
	return nil, nil
}

func getAwsRDSOptionGroupTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsRDSOptionGroupTags")

	optionGroup := h.Item.(types.OptionGroup)

	// Create service
	svc, err := RDSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_option_group.getAwsRDSOptionGroupTags", "connection_error", err)
		return nil, err
	}

	params := &rds.ListTagsForResourceInput{
		ResourceName: optionGroup.OptionGroupArn,
	}

	op, err := svc.ListTagsForResource(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_rds_db_option_group.getAwsRDSOptionGroupTags", "api_error", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func getRDSDBOptionGroupTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	optionGroup := d.HydrateItem.(*rds.ListTagsForResourceOutput)

	if optionGroup.TagList != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range optionGroup.TagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
