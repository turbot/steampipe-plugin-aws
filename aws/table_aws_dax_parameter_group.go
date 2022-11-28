package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dax"
	"github.com/aws/aws-sdk-go-v2/service/dax/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsDaxParameterGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_dax_parameter_group",
		Description: "AWS DAX Parameter Group",
		List: &plugin.ListConfig{
			Hydrate: listDaxParameterGroups,
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
				Name:        "parameter_group_name",
				Description: "The name of the parameter group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A description of the parameter group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "parameters",
				Description: "A list of parameters within a parameter group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDaxParameters,
				Transform:   transform.FromValue(),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ParameterGroupName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDaxParameterGroupsAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

type Parameter struct {
	AllowedValues          string
	ChangeType             string
	DataType               string
	Description            string
	IsModifiable           string
	NodeTypeSpecificValues []types.NodeTypeSpecificValue
	ParameterName          string
	ParameterType          string
	ParameterValue         string
	Source                 string
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
	equalQuals := d.KeyColumnQuals
	if equalQuals["parameter_group_name"] != nil {
		params.ParameterGroupNames = []string{equalQuals["parameter_group_name"].GetStringValue()}
	}

	for pagesLeft {
		result, err := svc.DescribeParameterGroups(ctx, params)
		if err != nil {
			logger.Error("aws_dax_parameter_group.listDaxParameterGroups", "api_error", err)
			return nil, err
		}

		for _, parameterGroup := range result.ParameterGroups {
			d.StreamListItem(ctx, parameterGroup)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
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

//// HYDRATE FUNCTION

func getDaxParameters(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	parameterGroup := h.Item.(types.ParameterGroup)

	// Create Client
	svc, err := DAXClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dax_cluster.getDaxParameters", "connection_error", err)
		return nil, err
	}

	params := &dax.DescribeParametersInput{
		ParameterGroupName: parameterGroup.ParameterGroupName,
	}

	var parameters []Parameter

	for {
		op, err := svc.DescribeParameters(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_dax_cluster.getDaxParameters", "api_error", err)
			return nil, err
		}

		for _, item := range op.Parameters {
			param := &Parameter{}
			if item.AllowedValues != nil {
				param.AllowedValues = *item.AllowedValues
			}
			if item.ChangeType != "" {
				param.ChangeType = string(item.ChangeType)
			}
			if item.DataType != nil {
				param.DataType = *item.DataType
			}
			if item.Description != nil {
				param.Description = *item.Description
			}
			if item.IsModifiable != "" {
				param.IsModifiable = string(item.IsModifiable)
			}
			if item.NodeTypeSpecificValues != nil {
				param.NodeTypeSpecificValues = item.NodeTypeSpecificValues
			}
			if item.ParameterName != nil {
				param.ParameterName = *item.ParameterName
			}
			if item.ParameterType != "" {
				param.ParameterType = string(item.ParameterType)
			}
			if item.ParameterValue != nil {
				param.ParameterValue = *item.ParameterValue
			}
			if item.Source != nil {
				param.Source = *item.Source
			}
			parameters = append(parameters, *param)
		}

		if op.NextToken != nil {
			params.NextToken = op.NextToken
		} else {
			break
		}
	}

	return parameters, nil
}

func getDaxParameterGroupsAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	name := *h.Item.(types.ParameterGroup).ParameterGroupName

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_dax_parameter_group.getDaxParameterGroupsAkas", "cache_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	akas := []string{"arn:" + commonColumnData.Partition + ":dax:" + region + "::parametergroup:" + name}

	return akas, nil
}
