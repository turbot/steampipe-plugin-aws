package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshiftserverless"
	"github.com/aws/aws-sdk-go-v2/service/redshiftserverless/types"

	redshiftserverlessv1 "github.com/aws/aws-sdk-go/service/redshiftserverless"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsRedshiftServerlessNamespace(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_redshiftserverless_namespace",
		Description: "AWS Redshift Serverless Namespace",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("namespace_name"),
			Hydrate:    getRedshiftServerlessNamespace,
		},
		List: &plugin.ListConfig{
			Hydrate: listRedshiftServerlessNamespaces,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(redshiftserverlessv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "namespace_name",
				Description: "The name of the namespace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "namespace_id",
				Description: "The unique identifier of the namespace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "namespace_arn",
				Description: "The Amazon Resource Name (ARN) that links to the namespace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the namespace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "admin_username",
				Description: "The username of the administrator for the first database created in the namespace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_date",
				Description: "The creation date of the namespace.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "db_name",
				Description: "The name of the first database created in the namespace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "default_iam_role_arn",
				Description: "The Amazon Resource Name (ARN) of the IAM role to set as a default in the namespace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kms_key_id",
				Description: "The ID of the Amazon Web Services Key Management Service key used to encrypt your data.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "iam_roles",
				Description: "A list of IAM roles to associate with the namespace.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "log_exports",
				Description: "The types of logs the namespace can export. Available export types are User log, Connection log, and User activity log.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "The list of tags for the namespace.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNamespaceTags,
				Transform:   transform.FromField("Tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("NamespaceName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNamespaceTags,
				Transform:   transform.From(getNamespaceTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("NamespaceArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listRedshiftServerlessNamespaces(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create Session
	svc, err := RedshiftServerlessClient(ctx, d)
	if err != nil {
		logger.Error("aws_redshiftserverless_namespace.listRedshiftServerlessNamespaces", "service_creation_error", err)
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

	input := &redshiftserverless.ListNamespacesInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := redshiftserverless.NewListNamespacesPaginator(svc, input, func(o *redshiftserverless.ListNamespacesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_redshiftserverless_namespace.listRedshiftServerlessNamespaces", "api_error", err)
			return nil, err
		}

		for _, namespace := range output.Namespaces {
			d.StreamListItem(ctx, namespace)
		}

		// Context may get cancelled due to manual cancellation or if the limit has been reached
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getRedshiftServerlessNamespace(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	name := d.EqualsQuals["namespace_name"].GetStringValue()
	// Return nil, if no input provided
	if name == "" {
		return nil, nil
	}

	// Create service
	svc, err := RedshiftServerlessClient(ctx, d)
	if err != nil {
		logger.Error("aws_redshiftserverless_namespace.getRedshiftServerlessNamespace", "service_creation_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build params
	params := &redshiftserverless.GetNamespaceInput{
		NamespaceName: aws.String(name),
	}

	op, err := svc.GetNamespace(ctx, params)
	if err != nil {
		logger.Error("aws_redshiftserverless_namespace.getRedshiftServerlessNamespace", "api_error", err)
		return nil, err
	}
	return *op.Namespace, nil
}

func getNamespaceTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	arn := *h.Item.(types.Namespace).NamespaceArn

	// Create service
	svc, err := RedshiftServerlessClient(ctx, d)
	if err != nil {
		logger.Error("aws_redshiftserverless_namespace.getNamespaceTags", "service_creation_error", err)
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
		logger.Error("aws_redshiftserverless_namespace.getNamespaceTags", "api_error", err)
		return nil, err
	}
	return op, nil
}

//// TRANSFORM FUNCTIONS

func getNamespaceTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
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
