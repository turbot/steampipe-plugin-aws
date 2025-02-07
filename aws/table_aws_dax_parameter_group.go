package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dax"

	daxEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsDaxParameterGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dax_parameter_group",
		Description: "AWS DAX Parameter Group",
		List: &plugin.ListConfig{
			Hydrate: listDaxParameterGroups,
			Tags:    map[string]string{"service": "dax", "action": "DescribeParameterGroups"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ParameterGroupNotFoundFault"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "parameter_group_name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(daxEndpoint.AWS_DAX_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "parameter_group_name",
				Description: "The name of the parameter group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "Description of the parameter group.",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ParameterGroupName"),
			},
		}),
	}
}

//// LIST FUNCTION

func listDaxParameterGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create Client
	svc, err := DAXClient(ctx, d)
	if err != nil {
		logger.Error("aws_dax_parameter_group.listDaxParameterGroups", "service_creation_error", err)
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
			if limit < 20 {
				maxLimit = 20
			} else {
				maxLimit = limit
			}
		}
	}

	pagesLeft := true
	params := &dax.DescribeParameterGroupsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	// Additonal Filter
	equalQuals := d.EqualsQuals
	if equalQuals["parameter_group_name"] != nil {
		params.ParameterGroupNames = []string{equalQuals["parameter_group_name"].GetStringValue()}
	}

	for pagesLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		result, err := svc.DescribeParameterGroups(ctx, params)
		if err != nil {
			logger.Error("aws_dax_parameter_group.listDaxParameterGroups", "api_error", err)
			return nil, err
		}

		for _, parameterGroup := range result.ParameterGroups {
			d.StreamListItem(ctx, parameterGroup)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if result.NextToken != nil {
			pagesLeft = true
			params.NextToken = result.NextToken
		} else {
			pagesLeft = false
		}
	}

	return nil, nil
}
