package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSecretsManagerSecret(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_secretsmanager_secret",
		Description: "AWS Secrets Manager Secret",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("arn"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ValidationException", "InvalidParameter", "ResourceNotFoundException"}),
			},
			Hydrate: describeSecretsManagerSecret,
		},
		List: &plugin.ListConfig{
			Hydrate: listSecretsManagerSecrets,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "name", Require: plugin.Optional},
				{Name: "description", Require: plugin.Optional},
				{Name: "primary_region", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
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
	svc, err := SecretsManagerService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &secretsmanager.ListSecretsInput{
		MaxResults: aws.Int64(100),
	}

	filters := buildSecretManagerSecretFilter(d.Quals)
	if len(filters) > 0 {
		input.Filters = filters
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			if *limit < 1 {
				input.MaxResults = aws.Int64(1)
			} else {
				input.MaxResults = limit
			}
		}
	}

	// List call
	err = svc.ListSecretsPages(
		input,
		func(page *secretsmanager.ListSecretsOutput, lastPage bool) bool {
			for _, secret := range page.SecretList {
				d.StreamListItem(ctx, secret)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !lastPage
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func describeSecretsManagerSecret(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("describeSecretsManagerSecret")

	var secretID string
	if h.Item != nil {
		data := secretData(h.Item)
		secretID = data["ARN"]
	} else {
		quals := d.KeyColumnQuals
		secretID = quals["arn"].GetStringValue()
	}

	// get service
	svc, err := SecretsManagerService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &secretsmanager.DescribeSecretInput{
		SecretId: aws.String(secretID),
	}

	// Get call
	op, err := svc.DescribeSecret(params)
	if err != nil {
		plugin.Logger(ctx).Debug("describeSecretsManagerSecret", "ERROR", err)
		return nil, err
	}
	return op, nil
}

func getSecretsManagerSecretPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getSecretsManagerSecretPolicy")

	var arn string
	if h.Item != nil {
		data := secretData(h.Item)
		arn = data["ARN"]
	}

	// Create Session
	svc, err := SecretsManagerService(ctx, d)
	if err != nil {
		logger.Error("getSecretsManagerSecretPolicy", "error_SecretsManagerService", err)
		return nil, err
	}

	// Build the params
	params := &secretsmanager.GetResourcePolicyInput{
		SecretId: &arn,
	}

	// Get call
	data, err := svc.GetResourcePolicy(params)
	if err != nil {
		logger.Error("getSecretsManagerSecretPolicy", "error_GetResourcePolicy", err)
		return nil, err
	}

	return data, nil
}

//// TRANSFORM FUNCTION

func secretsManagerSecretTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("secretsManagerSecretTagListToTurbotTags")
	tagList := d.Value.([]*secretsmanager.Tag)

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
	case *secretsmanager.SecretListEntry:
		data["ARN"] = *item.ARN
	}
	return data
}

//// UTILITY FUNCTION

// Build secret manager secret list call input filter
func buildSecretManagerSecretFilter(quals plugin.KeyColumnQualMap) []*secretsmanager.Filter {
	filters := make([]*secretsmanager.Filter, 0)

	filterQuals := map[string]string{
		"description":    "description",
		"name":           "name",
		"primary_region": "primary-region",
	}
	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := secretsmanager.Filter{
				Key: types.String(filterName),
			}
			value := getQualsValueByColumn(quals, columnName, "string")
			val, ok := value.(string)
			if ok {
				filter.Values = []*string{aws.String(val)}
			} else {
				filter.Values = value.([]*string)
			}
			filters = append(filters, &filter)
		}
	}
	return filters
}
