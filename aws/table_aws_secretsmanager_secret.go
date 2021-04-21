package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSecretsManagerSecret(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_secretsmanager_secret",
		Description: "AWS Secrets Manager Secret",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("arn, name"),
			Hydrate:           describeSecretsManagerSecret,
		},
		List: &plugin.ListConfig{
			Hydrate: listSecretsManagerSecrets,
		},
		GetMatrixItem: BuildRegionList,
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
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listSecretsManagerSecrets", "AWS_REGION", region)

	// Create session
	svc, err := SecretsManagerService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.ListSecretsPages(
		&secretsmanager.ListSecretsInput{},
		func(page *secretsmanager.ListSecretsOutput, lastPage bool) bool {
			for _, secret := range page.SecretList {
				d.StreamListItem(ctx, secret)
			}
			return !lastPage
		},
	)

	return nil, err
}

//// HYDRATE FUNCTION

func describeSecretsManagerSecret(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("describeSecretsManagerSecret")

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	var secretID string
	if h.Item != nil {
		data := secretData(h.Item)
		secretID = data["ARN"]
	} else {
		quals := d.KeyColumnQuals
		secretID = quals["arn"].GetStringValue()
	}

	// get service
	svc, err := SecretsManagerService(ctx, d, region)
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
