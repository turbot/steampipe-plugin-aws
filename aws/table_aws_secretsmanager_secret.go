package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"

	secretsmanagerEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSecretsManagerSecret(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_secretsmanager_secret",
		Description: "AWS Secrets Manager Secret",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationException", "InvalidParameter", "ResourceNotFoundException"}),
			},
			Hydrate: describeSecretsManagerSecret,
			Tags:    map[string]string{"service": "secretsmanager", "action": "DescribeSecret"},
		},
		List: &plugin.ListConfig{
			Hydrate: listSecretsManagerSecrets,
			Tags:    map[string]string{"service": "secretsmanager", "action": "ListSecrets"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "name", Require: plugin.Optional},
				{Name: "description", Require: plugin.Optional},
				{Name: "primary_region", Require: plugin.Optional},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: describeSecretsManagerSecret,
				Tags: map[string]string{"service": "sagemaker", "action": "DescribeSecret"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(secretsmanagerEndpoint.SECRETSMANAGERServiceID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name of the secret.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the secret.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ARN"),
			},
			{
				Name:        "created_date",
				Description: "The date and time when a secret was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "The user-provided description of the secret.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kms_key_id",
				Description: "The ARN or alias of the AWS KMS customer master key (CMK) used to encrypt the SecretString and SecretBinary fields in each version of the secret.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeSecretsManagerSecret,
			},
			{
				Name:        "deleted_date",
				Description: "The date and time the deletion of the secret occurred.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_accessed_date",
				Description: "The last date that this secret was accessed.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_changed_date",
				Description: "The last date and time that this secret was modified in any way.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_rotated_date",
				Description: "The most recent date and time that the Secrets Manager rotation process was successfully completed.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "owning_service",
				Description: "Returns the name of the service that created the secret.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeSecretsManagerSecret,
			},
			{
				Name:        "primary_region",
				Description: "The Region where Secrets Manager originated the secret.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeSecretsManagerSecret,
			},
			{
				Name:        "policy",
				Description: "A JSON-formatted string that describes the permissions that are associated with the attached secret.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSecretsManagerSecretPolicy,
				Transform:   transform.FromField("ResourcePolicy"),
			},
			{
				Name:        "policy_std",
				Description: "Contains the permissions that are associated with the attached secret in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSecretsManagerSecretPolicy,
				Transform:   transform.FromField("ResourcePolicy").Transform(unescape).Transform(policyToCanonical),
			},
			{
				Name:        "replication_status",
				Description: "Describes a list of replication status objects as InProgress, Failed or InSync.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeSecretsManagerSecret,
			},
			{
				Name:        "rotation_enabled",
				Description: "Indicates whether automatic, scheduled rotation is enabled for this secret.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "rotation_lambda_arn",
				Description: "The ARN of an AWS Lambda function invoked by Secrets Manager to rotate and expire the secret either automatically per the schedule or manually by a call to RotateSecret.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RotationLambdaARN"),
				Hydrate:     describeSecretsManagerSecret,
			},
			{
				Name:        "rotation_rules",
				Description: "A structure that defines the rotation configuration for the secret.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeSecretsManagerSecret,
			},
			{
				Name:        "secret_versions_to_stages",
				Description: "A list of all of the currently assigned SecretVersionStage staging labels and the SecretVersionId attached to each one.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "The list of user-defined tags associated with the secret.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeSecretsManagerSecret,
				Transform:   transform.FromField("Tags").Transform(secretsManagerSecretTagListToTurbotTags),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ARN").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listSecretsManagerSecrets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := SecretsManagerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_secretsmanager_secret.listSecretsManagerSecrets", "connection_error", err)
		return nil, err
	}

	// Limiting the results
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

	input := &secretsmanager.ListSecretsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	filters := buildSecretManagerSecretFilter(d.Quals)
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := secretsmanager.NewListSecretsPaginator(svc, input, func(o *secretsmanager.ListSecretsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_secretsmanager_secret.listSecretsManagerSecrets", "api_error", err)
			return nil, err
		}

		for _, items := range output.SecretList {
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

func describeSecretsManagerSecret(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var secretID string
	if h.Item != nil {
		data := secretData(h.Item)
		secretID = data["ARN"]
	} else {
		quals := d.EqualsQuals
		secretID = quals["arn"].GetStringValue()
	}

	// get service
	svc, err := SecretsManagerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_secretsmanager_secret.describeSecretsManagerSecret", "connection_error", err)
		return nil, err
	}

	// Build the params
	params := &secretsmanager.DescribeSecretInput{
		SecretId: aws.String(secretID),
	}

	// Get call
	op, err := svc.DescribeSecret(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_secretsmanager_secret.describeSecretsManagerSecret", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getSecretsManagerSecretPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var arn string
	if h.Item != nil {
		data := secretData(h.Item)
		arn = data["ARN"]
	}

	// Create Session
	svc, err := SecretsManagerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_secretsmanager_secret.getSecretsManagerSecretPolicy", "connection_eror", err)
		return nil, err
	}

	// Build the params
	params := &secretsmanager.GetResourcePolicyInput{
		SecretId: &arn,
	}

	// Get call
	data, err := svc.GetResourcePolicy(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_secretsmanager_secret.getSecretsManagerSecretPolicy", "api_error", err)
		return nil, err
	}

	return data, nil
}

//// TRANSFORM FUNCTION

func secretsManagerSecretTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	tagList := d.Value.([]types.Tag)

	// Mapping the resource tags inside turbotTags
	var turbotTagsMap map[string]string
	if tagList != nil {
		turbotTagsMap = map[string]string{}
		for _, i := range tagList {
			turbotTagsMap[*i.Key] = *i.Value
		}
	}

	return turbotTagsMap, nil
}

func secretData(item interface{}) map[string]string {
	data := map[string]string{}
	switch item := item.(type) {
	case *secretsmanager.DescribeSecretOutput:
		data["ARN"] = *item.ARN
	case types.SecretListEntry:
		data["ARN"] = *item.ARN
	}
	return data
}

//// UTILITY FUNCTION

// Build secret manager secret list call input filter
func buildSecretManagerSecretFilter(quals plugin.KeyColumnQualMap) []types.Filter {
	filters := make([]types.Filter, 0)

	filterQuals := map[string]string{
		"description":    "description",
		"name":           "name",
		"primary_region": "primary-region",
	}
	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := types.Filter{
				Key: types.FilterNameStringType(filterName),
			}
			value := getQualsValueByColumn(quals, columnName, "string")
			val, ok := value.(string)
			if ok {
				filter.Values = []string{val}
			}
			filters = append(filters, filter)
		}
	}
	return filters
}
