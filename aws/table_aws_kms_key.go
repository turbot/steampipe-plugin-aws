package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/kms"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsKmsKey(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_kms_key",
		Description: "AWS KMS Key",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getKmsKey,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"NotFoundException", "InvalidParameter"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listKmsKeys,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "Unique identifier of the key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("KeyId"),
			},
			{
				Name:        "arn",
				Description: "ARN of the key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("KeyArn"),
			},
			{
				Name:        "key_manager",
				Description: "The manager of the CMK. CMKs in your AWS account are either customer managed or AWS managed.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsKmsKeyData,
				Transform:   transform.FromField("KeyMetadata.KeyManager"),
			},
			{
				Name:        "creation_date",
				Description: "The date and time when the CMK was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getAwsKmsKeyData,
				Transform:   transform.FromField("KeyMetadata.CreationDate"),
			},
			{
				Name:        "aws_account_id",
				Description: "The twelve-digit account ID of the AWS account that owns the CMK.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsKmsKeyData,
				Transform:   transform.FromField("KeyMetadata.AWSAccountId"),
			},
			{
				Name:        "enabled",
				Description: "Specifies whether the CMK is enabled. When KeyState is Enabled this value is true, otherwise it is false.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsKmsKeyData,
				Transform:   transform.FromField("KeyMetadata.Enabled"),
			},
			{
				Name:        "customer_master_key_spec",
				Description: "Describes the type of key material in the CMK.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsKmsKeyData,
				Transform:   transform.FromField("KeyMetadata.CustomerMasterKeySpec"),
			},
			{
				Name:        "description",
				Description: "The description of the CMK.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsKmsKeyData,
				Transform:   transform.FromField("KeyMetadata.Description"),
			},
			{
				Name:        "deletion_date",
				Description: "The date and time after which AWS KMS deletes the CMK.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getAwsKmsKeyData,
				Transform:   transform.FromField("KeyMetadata.DeletionDate"),
			},
			{
				Name:        "key_state",
				Description: "The current status of the CMK. For more information about how key state affects the use of a CMK, see [Key state: Effect on your CMK](https://docs.aws.amazon.com/kms/latest/developerguide/key-state.html).",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsKmsKeyData,
				Transform:   transform.FromField("KeyMetadata.KeyState"),
			},
			{
				Name:        "key_usage",
				Description: "The [cryptographic operations](https://docs.aws.amazon.com/kms/latest/developerguide/concepts.html#cryptographic-operations) for which you can use the CMK.",
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
				Description: "The time at which the imported key material expires.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getAwsKmsKeyData,
				Transform:   transform.FromField("KeyMetadata.ValidTo"),
			},
			{
				Name:        "aliases",
				Description: "A list of aliases for the key.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsKmsKeyAliases,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "key_rotation_enabled",
				Description: "A Boolean value that specifies whether key rotation is enabled.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsKmsKeyRotationStatus,
			},
			{
				Name:        "policy",
				Description: "A key policy document in JSON format.",
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
				Description: "A list of tags attached to key.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsKmsKeyTagging,
			},

			/// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsKmsKeyAliases,
				Transform:   transform.From(kmsKeyTitle),
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

//// LIST FUNCTION

func listKmsKeys(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := KMSService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &kms.ListKeysInput{
		Limit: aws.Int64(1000),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.Limit {
			if *limit < 1 {
				input.Limit = aws.Int64(1)
			} else {
				input.Limit = limit
			}
		}
	}

	err = svc.ListKeysPages(
		input,
		func(page *kms.ListKeysOutput, lastPage bool) bool {
			for _, key := range page.Keys {
				d.StreamListItem(ctx, key)

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

func getKmsKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getKmsKey")

	keyID := d.KeyColumnQuals["id"].GetStringValue()

	// Create Session
	svc, err := KMSService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &kms.DescribeKeyInput{
		KeyId: aws.String(keyID),
	}

	keyData, err := svc.DescribeKey(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getIamUser__", "ERROR", err)
		return nil, err
	}

	rowData := &kms.KeyListEntry{
		KeyArn: keyData.KeyMetadata.Arn,
		KeyId:  keyData.KeyMetadata.KeyId,
	}

	return rowData, nil
}

func getAwsKmsKeyData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsKmsKeyData")
	key := h.Item.(*kms.KeyListEntry)

	// Create Session
	svc, err := KMSService(ctx, d)
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
	plugin.Logger(ctx).Trace("getAwsKmsKeyRotationStatus")
	key := h.Item.(*kms.KeyListEntry)

	// Create Session
	svc, err := KMSService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &kms.GetKeyRotationStatusInput{
		KeyId: key.KeyId,
	}

	keyData, err := svc.GetKeyRotationStatus(params)
	if err != nil {
		// For AWS managed KMS keys GetKeyRotationStatus API generates exceptions
		if a, ok := err.(awserr.Error); ok {
			if helpers.StringSliceContains([]string{"AccessDeniedException", "UnsupportedOperationException"}, a.Code()) {
				return kms.GetKeyRotationStatusOutput{}, nil
			}
		}
		return nil, err
	}
	return keyData, nil
}

func getAwsKmsKeyTagging(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsKmsKeyTagging")
	key := h.Item.(*kms.KeyListEntry)

	// Create Session
	svc, err := KMSService(ctx, d)
	if err != nil {
		return nil, err
	}

	tagsData := map[string]interface{}{}
	params := &kms.ListResourceTagsInput{
		KeyId: key.KeyId,
	}

	pagesLeft := true
	tags := []*kms.Tag{}
	for pagesLeft {
		keyTags, err := svc.ListResourceTags(params)
		if err != nil {
			// For AWS managed KMS keys ListResourceTags API generates AccessDeniedException
			if a, ok := err.(awserr.Error); ok {
				if a.Code() == "AccessDeniedException" {
					return tagsData, nil
				}
			}
			return nil, err
		}
		tags = append(tags, keyTags.Tags...)

		if keyTags.NextMarker != nil {
			params.Marker = keyTags.NextMarker
		} else {
			pagesLeft = false
		}
	}

	if tags != nil {
		tagsData["TagsSrc"] = tags

		turbotTagsMap := make(map[string]string)
		for _, t := range tags {
			turbotTagsMap[*t.TagKey] = *t.TagValue
		}
		tagsData["Tags"] = turbotTagsMap
	}

	return tagsData, nil
}

func getAwsKmsKeyAliases(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsKmsKeyAliases")
	key := h.Item.(*kms.KeyListEntry)

	// Create Session
	svc, err := KMSService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &kms.ListAliasesInput{
		KeyId: key.KeyId,
	}
	keyData := []*kms.AliasListEntry{}
	err = svc.ListAliasesPages(
		params,
		func(page *kms.ListAliasesOutput, lastPage bool) bool {
			keyData = append(keyData, page.Aliases...)
			return !lastPage
		},
	)
	if err != nil {
		plugin.Logger(ctx).Error("getAwsKmsKeyAliases", "ListAliasesPages_error", err)
	}

	return keyData, nil
}

func kmsKeyTitle(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	// Use the first alias if one is set, else fallback to the key ID
	key := d.HydrateItem.([]*kms.AliasListEntry)
	if len(key) > 0 {
		return key[0].AliasName, nil
	}

	var keyID string
	if d.HydrateResults["listKmsKeys"] != nil {
		keyID = *(d.HydrateResults["listKmsKeys"]).(*kms.KeyListEntry).KeyId
	} else if d.HydrateResults["getKmsKey"] != nil {
		keyID = *(d.HydrateResults["getKmsKey"]).(*kms.KeyListEntry).KeyId
	}

	return keyID, nil
}

func getAwsKmsKeyPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsKmsKeyPolicy")
	key := h.Item.(*kms.KeyListEntry)

	// Create Session
	svc, err := KMSService(ctx, d)
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
