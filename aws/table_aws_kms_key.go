package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/kms"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsKmsKey(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_kms_key",
		Description: "AWS KMS Key",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"NotFoundException"}),
			ItemFromKey:       kmsKeyFromKey,
			Hydrate:           getKmsKey,
		},
		List: &plugin.ListConfig{
			Hydrate: listKmsKeys,
		},
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "Unique identifier of the key",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("KeyId"),
			},
			{
				Name:        "arn",
				Description: "ARN of the key",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("KeyArn"),
			},
			{
				Name:        "key_manager",
				Description: "The manager of the CMK. CMKs in your AWS account are either customer managed or AWS managed",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsKmsKeyData,
				Transform:   transform.FromField("KeyMetadata.KeyManager"),
			},
			{
				Name:        "creation_date",
				Description: "The date and time when the CMK was created",
				Type:        proto.ColumnType_DATETIME,
				Hydrate:     getAwsKmsKeyData,
				Transform:   transform.FromField("KeyMetadata.CreationDate"),
			},
			{
				Name:        "aws_account_id",
				Description: "The twelve-digit account ID of the AWS account that owns the CMK",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsKmsKeyData,
				Transform:   transform.FromField("KeyMetadata.AWSAccountId"),
			},
			{
				Name:        "enabled",
				Description: "Specifies whether the CMK is enabled. When KeyState is Enabled this value is true, otherwise it is false",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsKmsKeyData,
				Transform:   transform.FromField("KeyMetadata.Enabled"),
			},
			{
				Name:        "customer_master_key_spec",
				Description: "Describes the type of key material in the CMK",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsKmsKeyData,
				Transform:   transform.FromField("KeyMetadata.CustomerMasterKeySpec"),
			},
			{
				Name:        "description",
				Description: "The description of the CMK",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsKmsKeyData,
				Transform:   transform.FromField("KeyMetadata.Description"),
			},
			{
				Name:        "deletion_date",
				Description: "The date and time after which AWS KMS deletes the CMK",
				Type:        proto.ColumnType_DATETIME,
				Hydrate:     getAwsKmsKeyData,
				Transform:   transform.FromField("KeyMetadata.DeletionDate"),
			},
			{
				Name:        "key_state",
				Description: "The current status of the CMK. For more information about how key state affects the use of a CMK, see [Key state: Effect on your CMK](https://docs.aws.amazon.com/kms/latest/developerguide/key-state.html)",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsKmsKeyData,
				Transform:   transform.FromField("KeyMetadata.KeyState"),
			},
			{
				Name:        "key_usage",
				Description: "The [cryptographic operations](https://docs.aws.amazon.com/kms/latest/developerguide/concepts.html#cryptographic-operations) for which you can use the CMK",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsKmsKeyData,
				Transform:   transform.FromField("KeyMetadata.KeyUsage"),
			},
			{
				Name:        "origin",
				Description: "The source of the CMK's key material. When this value is AWS_KMS, AWS KMS created the key material. When this value is EXTERNAL, the key material was imported from your existing key management infrastructure or the CMK lacks key material. When this value is AWS_CLOUDHSM, the key material was created in the AWS CloudHSM cluster associated with a custom key store.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsKmsKeyData,
				Transform:   transform.FromField("KeyMetadata.Origin"),
			},
			{
				Name:        "valid_to",
				Description: "The time at which the imported key material expires",
				Type:        proto.ColumnType_DATETIME,
				Hydrate:     getAwsKmsKeyData,
				Transform:   transform.FromField("KeyMetadata.ValidTo"),
			},
			{
				Name:        "aliases",
				Description: "A list of aliases for the key",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsKmsKeyAliases,
			},
			{
				Name:        "key_rotation_enabled",
				Description: "A Boolean value that specifies whether key rotation is enabled",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsKmsKeyRotationStatus,
			},
			{
				Name:        "policy",
				Description: "A key policy document in JSON format",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsKmsKeyPolicy,
				Transform:   transform.FromField("Policy").Transform(transform.UnmarshalYAML),
			},
			{
				Name:        "policy_std",
				Description: "Contains the policy in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsKmsKeyPolicy,
				Transform:   transform.FromField("Policy").Transform(unescape).Transform(policyToCanonical),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached to key",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsKmsKeyTagging,
			},

			/// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("KeyId"),
				// Default:     transform.FromField("KeyId"),
				// Hydrate:     getAwsKmsKeyAliases,
				// Transform:   transform.From(getAwsKmsKeyTitle),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsKmsKeyTagging,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("KeyArn").Transform(arnToAkas),
			},
		}),
	}
}

//// BUILD HYDRATE INPUT

func kmsKeyFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	ID := quals["id"].GetStringValue()
	item := &kms.KeyListEntry{
		KeyId: &ID,
	}
	return item, nil
}

//// LIST FUNCTION

func listKmsKeys(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()

	// Create Session
	svc, err := KMSService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	err = svc.ListKeysPages(
		&kms.ListKeysInput{},
		func(page *kms.ListKeysOutput, lastPage bool) bool {
			for _, key := range page.Keys {
				d.StreamListItem(ctx, key)
			}
			return true
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getKmsKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getKmsKey")
	key := h.Item.(*kms.KeyListEntry)
	defaultRegion := GetDefaultRegion()

	// Create Session
	svc, err := KMSService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &kms.DescribeKeyInput{
		KeyId: key.KeyId,
	}

	keyData, err := svc.DescribeKey(params)
	if err != nil {
		logger.Debug("getIamUser__", "ERROR", err)
		return nil, err
	}

	var rowData *kms.KeyListEntry
	rowData = &kms.KeyListEntry{
		KeyArn: keyData.KeyMetadata.Arn,
		KeyId:  keyData.KeyMetadata.KeyId,
	}

	return rowData, nil
}

func getAwsKmsKeyData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsKmsKeyData")
	key := h.Item.(*kms.KeyListEntry)
	defaultRegion := GetDefaultRegion()

	// Create Session
	svc, err := KMSService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &kms.DescribeKeyInput{
		KeyId: key.KeyId,
	}

	keyData, err := svc.DescribeKey(params)
	if err != nil {
		return nil, err
	}

	return keyData, nil
}

func getAwsKmsKeyRotationStatus(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsKmsKeyRotationStatus")
	key := h.Item.(*kms.KeyListEntry)
	defaultRegion := GetDefaultRegion()

	// Create Session
	svc, err := KMSService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &kms.GetKeyRotationStatusInput{
		KeyId: key.KeyId,
	}

	keyData, err := svc.GetKeyRotationStatus(params)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "AccessDeniedException" {
				return kms.GetKeyRotationStatusOutput{}, nil
			}
		}
		return nil, err
	}
	return keyData, nil
}

func getAwsKmsKeyTagging(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsKmsKeyTagging")
	key := h.Item.(*kms.KeyListEntry)
	defaultRegion := GetDefaultRegion()

	// Create Session
	svc, err := KMSService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	tagsData := map[string]interface{}{}
	params := &kms.ListResourceTagsInput{
		KeyId: key.KeyId,
	}

	keyTags, err := svc.ListResourceTags(params)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "AccessDeniedException" {
				return tagsData, nil
			}
		}
		return nil, err
	}
	if keyTags.Tags != nil {
		tagsData["TagsSrc"] = keyTags.Tags

		turbotTagsMap := make(map[string]string)
		for _, t := range keyTags.Tags {
			turbotTagsMap[*t.TagKey] = *t.TagValue
		}
		tagsData["Tags"] = turbotTagsMap
	}
	return tagsData, nil
}

func getAwsKmsKeyAliases(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsKmsKeyAliases")
	key := h.Item.(*kms.KeyListEntry)
	defaultRegion := GetDefaultRegion()

	// Create Session
	svc, err := KMSService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &kms.ListAliasesInput{
		KeyId: key.KeyId,
	}

	keyData, err := svc.ListAliases(params)
	if err != nil {
		return nil, err
	}

	return keyData, nil
}

func getAwsKmsKeyPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getAwsKmsKeyPolicy")
	key := h.Item.(*kms.KeyListEntry)
	defaultRegion := GetDefaultRegion()

	// Create Session
	svc, err := KMSService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &kms.GetKeyPolicyInput{
		KeyId:      key.KeyId,
		PolicyName: aws.String("default"),
	}

	keyPolicy, err := svc.GetKeyPolicy(params)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "NotFoundException" {
				return kms.GetKeyPolicyOutput{}, nil
			}
		}
		return nil, err
	}
	return keyPolicy, nil
}

//// TRANSFORM FUNCTIONS//

func getAwsKmsKeyTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	hydrateData := d.HydrateItem.(*kms.ListAliasesOutput)
	if hydrateData.Aliases != nil && len(hydrateData.Aliases) > 0 {
		return hydrateData.Aliases[0].AliasName, nil
	}
	return nil, nil
}
