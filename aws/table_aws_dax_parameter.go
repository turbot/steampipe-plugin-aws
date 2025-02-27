package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dax"
	"github.com/aws/aws-sdk-go-v2/service/dax/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsDaxParameter(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dax_parameter",
		Description: "AWS DAX Parameter",
		List: &plugin.ListConfig{
			ParentHydrate: listDaxParameterGroups,
			Hydrate:       listDaxParameters,
			Tags:          map[string]string{"service": "dax", "action": "DescribeParameters"},
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
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_DAX_SERVICE_ID),
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
				Description: "The value of the parameter.",
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
				Description: "The conditions under which changes to this parameter can be applied. Possible values are 'IMMEDIATE', 'REQUIRES_REBOOT'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "data_type",
				Description: "The data type of the parameter.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_modifiable",
				Description: "Whether the customer is allowed to modify the parameter. Possible values are 'TRUE', 'FALSE' 'CONDITIONAL'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "parameter_type",
				Description: "Determines whether the parameter can be applied to any node or only nodes of a particular type. Possible values are 'DEFAULT', 'NODE_TYPE_SPECIFIC'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source",
				Description: "How the parameter is defined. For example, system denotes a system-defined parameter.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "node_type_specific_values",
				Description: "A list of node types, and specific parameter values for each node.",
				Type:        proto.ColumnType_JSON,
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
	equalQuals := d.EqualsQuals
	if equalQuals["parameter_group_name"] != nil {
		if *parameterGroup.ParameterGroupName != d.EqualsQuals["parameter_group_name"].GetStringValue() {
			return nil, nil
		}
	}

	// Create Client
	svc, err := DAXClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dax_cluster. listDaxParameters", "connection_error", err)
		return nil, err
	}

	params := &dax.DescribeParametersInput{
		ParameterGroupName: parameterGroup.ParameterGroupName,
	}

	for {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		op, err := svc.DescribeParameters(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_dax_cluster. listDaxParameters", "api_error", err)
			return nil, err
		}

		for _, item := range op.Parameters {
			d.StreamListItem(ctx, &Parameter{*parameterGroup.ParameterGroupName, item})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
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
