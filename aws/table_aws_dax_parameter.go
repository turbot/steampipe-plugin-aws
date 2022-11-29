package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go-v2/service/dax"
	"github.com/aws/aws-sdk-go-v2/service/dax/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsDaxParameter(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dax_parameter",
		Description: "AWS DAX Parameter",
		List: &plugin.ListConfig{
			ParentHydrate: listDaxParameterGroups,
			Hydrate:       listDaxParameters,
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
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "parameter_name",
				Description: "The name of the parameter.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "parameter_group_name",
				Description: "The name of the parameter group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "parameter_value",
				Description: "The value for the parameter.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "Description of the parameter.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "allowed_values",
				Description: "A range of values within which the parameter can be set.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "change_type",
				Description: "The conditions under which changes to this parameter can be applied.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "data_type",
				Description: "The data type of the parameter.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_modifiable",
				Description: "Whether the customer is allowed to modify the parameter.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "parameter_type",
				Description: "Determines whether the parameter can be applied to any nodes, or only nodes of a particular type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source",
				Description: "How the parameter is defined. For example, system denotes a system-defined parameter.",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ParameterName"),
			},
		}),
	}
}

type Parameter struct {
	ParameterGroupName string
	types.Parameter
}

//// LIST FUNCTION

func listDaxParameters(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	parameterGroup := h.Item.(types.ParameterGroup)

	// Additonal Filter
	equalQuals := d.KeyColumnQuals
	if equalQuals["parameter_group_name"] != nil {
		if *parameterGroup.ParameterGroupName != d.KeyColumnQuals["parameter_group_name"].GetStringValue() {
			return nil, nil
		}
	}

	// Create Client
	svc, err := DAXClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dax_cluster.getDaxParameters", "connection_error", err)
		return nil, err
	}

	params := &dax.DescribeParametersInput{
		ParameterGroupName: parameterGroup.ParameterGroupName,
	}

	for {
		op, err := svc.DescribeParameters(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_dax_cluster.getDaxParameters", "api_error", err)
			return nil, err
		}

		for _, item := range op.Parameters {
			d.StreamListItem(ctx, &Parameter{*parameterGroup.ParameterGroupName, item})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if op.NextToken != nil {
			params.NextToken = op.NextToken
		} else {
			break
		}
	}

	return nil, nil
}
