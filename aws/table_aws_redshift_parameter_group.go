package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
	"github.com/aws/aws-sdk-go-v2/service/redshift/types"

	redshiftEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsRedshiftParameterGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_redshift_parameter_group",
		Description: "AWS Redshift Parameter Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ClusterParameterGroupNotFound"}),
			},
			Hydrate: getRedshiftParameterGroup,
			Tags:    map[string]string{"service": "redshift", "action": "DescribeClusterParameterGroups"},
		},
		List: &plugin.ListConfig{
			Hydrate: listRedshiftParameterGroups,
			Tags:    map[string]string{"service": "redshift", "action": "DescribeClusterParameterGroups"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getRedshiftParameters,
				Tags: map[string]string{"service": "redshift", "action": "DescribeClusterParameters"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(redshiftEndpoint.REDSHIFTServiceID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the cluster parameter group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ParameterGroupName"),
			},
			{
				Name:        "description",
				Description: "The description of the parameter group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "family",
				Description: "The name of the cluster parameter group family that this cluster parameter group is compatible with.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ParameterGroupFamily"),
			},
			{
				Name:        "parameters",
				Description: "A list of Parameter instances. Each instance lists the parameters of one cluster parameter group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRedshiftParameters,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the parameter group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ParameterGroupName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags").Transform(tagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsRedshiftParameterGroupAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listRedshiftParameterGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := RedshiftClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_redshift_parameter_group.listRedshiftParameterGroups", "connection_error", err)
		return nil, err
	}

	input := &redshift.DescribeClusterParameterGroupsInput{
		MaxRecords: aws.Int32(100),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < *input.MaxRecords {
			if limit < 20 {
				input.MaxRecords = aws.Int32(20)
			} else {
				input.MaxRecords = aws.Int32(limit)
			}
		}
	}

	// List call
	paginator := redshift.NewDescribeClusterParameterGroupsPaginator(svc, input, func(o *redshift.DescribeClusterParameterGroupsPaginatorOptions) {
		o.Limit = *input.MaxRecords
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_redshift_parameter_group.listRedshiftParameterGroups", "api_error", err)
			return nil, err
		}

		for _, items := range output.ParameterGroups {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getRedshiftParameterGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := RedshiftClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_redshift_parameter_group.getRedshiftParameterGroup", "connection_error", err)
		return nil, err
	}

	name := d.EqualsQuals["name"].GetStringValue()

	// Return nil, if no input provided
	if name == "" {
		return nil, nil
	}

	// Build the params
	params := &redshift.DescribeClusterParameterGroupsInput{
		ParameterGroupName: aws.String(name),
	}

	// Get call
	data, err := svc.DescribeClusterParameterGroups(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_redshift_parameter_group.getRedshiftParameterGroup", "api_error", err)
		return nil, err
	}

	if len(data.ParameterGroups) > 0 {
		return data.ParameterGroups[0], nil
	}

	return nil, nil
}

func getRedshiftParameters(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := RedshiftClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_redshift_parameter_group.getRedshiftParameters", "connection_error", err)
		return nil, err
	}

	name := h.Item.(types.ClusterParameterGroup).ParameterGroupName

	// Build the params
	params := &redshift.DescribeClusterParametersInput{
		ParameterGroupName: name,
	}

	// Get call
	op, err := svc.DescribeClusterParameters(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_redshift_parameter_group.getRedshiftParameters", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getAwsRedshiftParameterGroupAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	parameterData := h.Item.(types.ClusterParameterGroup)

	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_redshift_parameter_group.getAwsRedshiftParameterGroupAkas", "getCommonColumns_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	aka := "arn:" + commonColumnData.Partition + ":redshift:" + region + ":" + commonColumnData.AccountId + ":parametergroup"

	if strings.HasPrefix(*parameterData.ParameterGroupName, ":") {
		aka = aka + *parameterData.ParameterGroupName
	} else {
		aka = aka + ":" + *parameterData.ParameterGroupName
	}

	return []string{aka}, nil
}

//// TRANSFORM FUNCTIONS

func tagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tagList := d.HydrateItem.(types.ClusterParameterGroup)

	if len(tagList.Tags) > 0 {
		turbotTagsMap := map[string]string{}
		for _, i := range tagList.Tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}
	return nil, nil
}
