package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/aws/aws-sdk-go-v2/service/glue/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsGlueMLTransform(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_glue_ml_transform",
		Description: "AWS Glue ML Transform",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("transform_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"EntityNotFoundException"}),
			},
			Hydrate: getGlueMLTransform,
			Tags:    map[string]string{"service": "glue", "action": "GetMLTransform"},
		},
		List: &plugin.ListConfig{
			Hydrate: listGlueMLTransforms,
			Tags:    map[string]string{"service": "glue", "action": "GetMLTransforms"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationException"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
				{
					Name:    "status",
					Require: plugin.Optional,
				},
				{
					Name:    "transform_type",
					Require: plugin.Optional,
				},
				{
					Name:    "glue_version",
					Require: plugin.Optional,
				},
				{
					Name:      "created_on",
					Require:   plugin.Optional,
					Operators: []string{">", ">=", "<", "<="},
				},
				{
					Name:      "last_modified_on",
					Require:   plugin.Optional,
					Operators: []string{">", ">=", "<", "<="},
				},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getGlueMLTransformTags,
				Tags: map[string]string{"service": "glue", "action": "GetTags"},
			},
			{
				Func: getGlueMLTransform,
				Tags: map[string]string{"service": "glue", "action": "GetMLTransform"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_GLUE_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "transform_id",
				Description: "The unique identifier of the transform, generated at the time that the transform was created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "The unique name given to the transform when it was created.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getGlueMLTransform,
			},
			{
				Name:        "description",
				Description: "A description of the transform.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getGlueMLTransform,
			},
			{
				Name:        "status",
				Description: "The last known status of the transform (to indicate whether it can be used or not). One of 'NOT_READY', 'READY', or 'DELETING'.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getGlueMLTransform,
			},
			{
				Name:        "created_on",
				Description: "The date and time when the transform was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getGlueMLTransform,
			},
			{
				Name:        "last_modified_on",
				Description: "The date and time when the transform was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getGlueMLTransform,
			},
			{
				Name:        "glue_version",
				Description: "This value determines which version of Glue this machine learning transform is compatible with.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getGlueMLTransform,
			},
			{
				Name:        "label_count",
				Description: "The number of labels available for this transform.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getGlueMLTransform,
			},
			{
				Name:        "max_capacity",
				Description: "The number of Glue data processing units (DPUs) that are allocated to task runs for this transform.",
				Type:        proto.ColumnType_DOUBLE,
				Hydrate:     getGlueMLTransform,
			},
			{
				Name:        "max_retries",
				Description: "The maximum number of times to retry a task for this transform after a task run fails.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getGlueMLTransform,
			},
			{
				Name:        "number_of_workers",
				Description: "The number of workers of a defined workerType that are allocated when this task runs.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getGlueMLTransform,
			},
			{
				Name:        "worker_type",
				Description: "The type of predefined worker that is allocated when this task runs. Accepts a value of Standard, G.1X, or G.2X.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getGlueMLTransform,
			},
			{
				Name:        "transform_type",
				Description: "The type of machine learning transform.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getGlueMLTransform,
				Transform:   transform.FromField("Parameters.TransformType"),
			},
			{
				Name:        "timeout",
				Description: "The timeout for a task run for this transform in minutes.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getGlueMLTransform,
			},
			{
				Name:        "role",
				Description: "The name or Amazon Resource Name (ARN) of the IAM role with the required permissions.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getGlueMLTransform,
			},
			{
				Name:        "input_record_tables",
				Description: "A list of Glue table definitions used by the transform.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGlueMLTransform,
			},
			{
				Name:        "parameters",
				Description: "The configuration parameters that are specific to the algorithm used.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGlueMLTransform,
			},
			{
				Name:        "schema",
				Description: "The Map object that represents the schema that this transform accepts.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGlueMLTransform,
			},
			{
				Name:        "evaluation_metrics",
				Description: "The latest evaluation metrics.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGlueMLTransform,
			},
			{
				Name:        "transform_encryption",
				Description: "The encryption-at-rest settings of the transform that apply to accessing user data.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGlueMLTransform,
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: "A map of tags for the resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGlueMLTransformTags,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Hydrate:     getGlueMLTransform,
				Transform:   transform.From(getGlueMLTransformTitle),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGlueMLTransformArn,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listGlueMLTransforms(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service
	svc, err := GlueClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_ml_transform.listGlueMLTransforms", "connection_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &glue.ListMLTransformsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	// Build filter criteria based on qualifiers
	filterCriteria := &types.TransformFilterCriteria{}

	if d.EqualsQuals["name"] != nil {
		filterCriteria.Name = aws.String(d.EqualsQuals["name"].GetStringValue())
	}

	if d.EqualsQuals["status"] != nil {
		status := types.TransformStatusType(d.EqualsQuals["status"].GetStringValue())
		filterCriteria.Status = status
	}

	if d.EqualsQuals["transform_type"] != nil {
		transformType := types.TransformType(d.EqualsQuals["transform_type"].GetStringValue())
		filterCriteria.TransformType = transformType
	}

	if d.EqualsQuals["glue_version"] != nil {
		filterCriteria.GlueVersion = aws.String(d.EqualsQuals["glue_version"].GetStringValue())
	}

	// Handle created_on qualifiers with operators
	quals := d.Quals
	if quals["created_on"] != nil {
		for _, q := range quals["created_on"].Quals {
			createdTime := q.Value.GetTimestampValue().AsTime()
			switch q.Operator {
			case ">=", ">":
				filterCriteria.CreatedAfter = &createdTime
			case "<", "<=":
				filterCriteria.CreatedBefore = &createdTime
			}
		}
	}

	// Handle last_modified_on qualifiers with operators
	if quals["last_modified_on"] != nil {
		for _, q := range quals["last_modified_on"].Quals {
			modifiedTime := q.Value.GetTimestampValue().AsTime()
			switch q.Operator {
			case ">=", ">":
				filterCriteria.LastModifiedAfter = &modifiedTime
			case "<", "<=":
				filterCriteria.LastModifiedBefore = &modifiedTime
			}
		}
	}

	// Only add filter if any criteria are set
	if filterCriteria.Name != nil || filterCriteria.Status != "" || filterCriteria.TransformType != "" ||
		filterCriteria.GlueVersion != nil || filterCriteria.CreatedAfter != nil || filterCriteria.CreatedBefore != nil ||
		filterCriteria.LastModifiedAfter != nil || filterCriteria.LastModifiedBefore != nil {
		input.Filter = filterCriteria
	}

	paginator := glue.NewListMLTransformsPaginator(svc, input, func(o *glue.ListMLTransformsPaginatorOptions) {
		o.Limit = maxLimit
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_glue_ml_transform.listGlueMLTransforms", "api_error", err)
			return nil, err
		}

		for _, transformId := range output.TransformIds {
			
			transformOutput := &glue.GetMLTransformOutput{
				TransformId: aws.String(transformId),
			}

			d.StreamListItem(ctx, transformOutput)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getGlueMLTransform(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var transformId string

	if h.Item != nil {
		transform := h.Item.(*glue.GetMLTransformOutput)
		transformId = *transform.TransformId
	} else {
		transformId = d.EqualsQuals["transform_id"].GetStringValue()
	}

	// Create service
	svc, err := GlueClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_ml_transform.getGlueMLTransform", "connection_error", err)
		return nil, err
	}

	params := &glue.GetMLTransformInput{
		TransformId: aws.String(transformId),
	}

	op, err := svc.GetMLTransform(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_ml_transform.getGlueMLTransform", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getGlueMLTransformTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the transform details
	var transformId string
	if h.Item != nil {
		transform := h.Item.(*glue.GetMLTransformOutput)
		transformId = *transform.TransformId
	}

	// Create service
	svc, err := GlueClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_ml_transform.getGlueMLTransformTags", "connection_error", err)
		return nil, err
	}

	// Get common columns
	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_ml_transform.getGlueMLTransformTags", "common_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)

	// Get the ARN for the transform
	arn := "arn:" + commonColumnData.Partition + ":glue:" + d.EqualsQualString(matrixKeyRegion) + ":" + commonColumnData.AccountId + ":mlTransform/" + transformId

	params := &glue.GetTagsInput{
		ResourceArn: aws.String(arn),
	}

	op, err := svc.GetTags(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_ml_transform.getGlueMLTransformTags", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getGlueMLTransformArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)

	var transformId string
	if h.Item != nil {
		transform := h.Item.(*glue.GetMLTransformOutput)
		transformId = *transform.TransformId
	} else {
		transformId = d.EqualsQuals["transform_id"].GetStringValue()
	}

	// Get common columns
	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glue_ml_transform.getGlueMLTransformArn", "common_error", err)
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)

	arn := "arn:" + commonColumnData.Partition + ":glue:" + region + ":" + commonColumnData.AccountId + ":mlTransform/" + transformId

	return arn, nil
}

//// TRANSFORM FUNCTIONS

func getGlueMLTransformTitle(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*glue.GetMLTransformOutput)

	if data.Name != nil {
		return *data.Name, nil
	}

	return nil, nil
}
