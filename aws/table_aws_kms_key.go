package aws

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"

	kmsv1 "github.com/aws/aws-sdk-go/service/kms"

	"github.com/aws/smithy-go"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsKmsKey(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_kms_key",
		Description: "AWS KMS Key",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getKmsKey,
			Tags:       map[string]string{"service": "kms", "action": "DescribeKey"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException", "InvalidParameter"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listKmsKeys,
			Tags:    map[string]string{"service": "kms", "action": "ListKeys"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(kmsv1.EndpointsID),
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsKmsKeyData,
				Tags: map[string]string{"service": "kms", "action": "DescribeKey"},
			},
			{
				Func: getAwsKmsKeyAliases,
				Tags: map[string]string{"service": "kms", "action": "ListAliases"},
			},
			{
				Func: getAwsKmsKeyRotationStatus,
				Tags: map[string]string{"service": "kms", "action": "GetKeyRotationStatus"},
			},
			{
				Func: getAwsKmsKeyPolicy,
				Tags: map[string]string{"service": "kms", "action": "GetKeyPolicy"},
			},
			{
				Func: getAwsKmsKeyTagging,
				Tags: map[string]string{"service": "kms", "action": "ListResourceTags"},
			},
		},
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
			{
				Name:        "multi_region",
				Description: "Specifies whether the CMK is KMS key is a multi-Region (true) or regional (false) key.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsKmsKeyData,
				Transform:   transform.FromField("KeyMetadata.MultiRegion"),
			},
			{
				Name:        "multi_region_configuration",
				Description: "Lists the primary and replica keys in same multi-Region key.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsKmsKeyData,
				Transform:   transform.FromField("KeyMetadata.MultiRegionConfiguration"),
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
	svc, err := KMSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kms_key.listKmsKeys", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	maxItems := int32(1000)
	input := &kms.ListKeysInput{}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
		}
	}
	input.Limit = aws.Int32(maxItems)
	paginator := kms.NewListKeysPaginator(svc, input, func(o *kms.ListKeysPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_kms_key.listKmsKeys", "api_error", err)
			return nil, err
		}

		for _, key := range output.Keys {
			d.StreamListItem(ctx, key)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getKmsKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	keyID := d.EqualsQuals["id"].GetStringValue()

	// Empty id check
	if strings.TrimSpace(keyID) == "" {
		return nil, nil
	}

	// Create Session
	svc, err := KMSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kms_key.getKmsKey", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &kms.DescribeKeyInput{
		KeyId: aws.String(keyID),
	}

	keyData, err := svc.DescribeKey(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kms_key.getKmsKey", "api_error", err)
		return nil, err
	}

	rowData := types.KeyListEntry{
		KeyArn: keyData.KeyMetadata.Arn,
		KeyId:  keyData.KeyMetadata.KeyId,
	}

	return rowData, nil
}

func getAwsKmsKeyData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	key := h.Item.(types.KeyListEntry)

	// Create Session
	svc, err := KMSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kms_key.getAwsKmsKeyData", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &kms.DescribeKeyInput{
		KeyId: key.KeyId,
	}

	keyData, err := svc.DescribeKey(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kms_key.getAwsKmsKeyData", "api_error", err)
		return nil, err
	}

	return keyData, nil
}

func getAwsKmsKeyRotationStatus(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	key := h.Item.(types.KeyListEntry)

	// Create Session
	svc, err := KMSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kms_key.getAwsKmsKeyRotationStatus", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &kms.GetKeyRotationStatusInput{
		KeyId: key.KeyId,
	}

	keyData, err := svc.GetKeyRotationStatus(ctx, params)
	if err != nil {
		// For AWS managed KMS keys GetKeyRotationStatus API generates exceptions
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if helpers.StringSliceContains([]string{"AccessDeniedException", "UnsupportedOperationException"}, ae.ErrorCode()) {
				return &kms.GetKeyRotationStatusOutput{}, nil
			}
			return nil, err
		}
		plugin.Logger(ctx).Error("aws_kms_key.getAwsKmsKeyRotationStatus", "api_error", err)
	}
	return keyData, nil
}

func getAwsKmsKeyTagging(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	key := h.Item.(types.KeyListEntry)

	// Create Session
	svc, err := KMSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kms_key.getAwsKmsKeyTagging", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	tagsData := map[string]interface{}{}
	params := &kms.ListResourceTagsInput{
		KeyId: key.KeyId,
	}

	tags := []types.Tag{}

	paginator := kms.NewListResourceTagsPaginator(svc, params, func(o *kms.ListResourceTagsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_kms_key.getAwsKmsKeyTagging", "api_error", err)
			return nil, err
		}
		tags = append(tags, output.Tags...)
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
	key := h.Item.(types.KeyListEntry)

	// Create Session
	svc, err := KMSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kms_key.getAwsKmsKeyAliases", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &kms.ListAliasesInput{
		KeyId: key.KeyId,
	}

	keyData := []types.AliasListEntry{}
	paginator := kms.NewListAliasesPaginator(svc, params, func(o *kms.ListAliasesPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_kms_key.getAwsKmsKeyAliases", "api_error", err)
			return nil, err
		}
		keyData = append(keyData, output.Aliases...)
	}

	return keyData, nil
}

func kmsKeyTitle(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	// Use the first alias if one is set, else fallback to the key ID
	key := d.HydrateItem.([]types.AliasListEntry)
	if len(key) > 0 {
		return key[0].AliasName, nil
	}

	var keyID string
	if d.HydrateResults["listKmsKeys"] != nil {
		keyID = *(d.HydrateResults["listKmsKeys"]).(types.KeyListEntry).KeyId
	} else if d.HydrateResults["getKmsKey"] != nil {
		keyID = *(d.HydrateResults["getKmsKey"]).(types.KeyListEntry).KeyId
	}

	return keyID, nil
}

func getAwsKmsKeyPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	key := h.Item.(types.KeyListEntry)

	// Create Session
	svc, err := KMSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kms_key.getAwsKmsKeyPolicy", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &kms.GetKeyPolicyInput{
		KeyId:      key.KeyId,
		PolicyName: aws.String("default"),
	}

	keyPolicy, err := svc.GetKeyPolicy(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "NotFoundException" {
				return &kms.GetKeyPolicyOutput{}, nil
			}
		}
		plugin.Logger(ctx).Error("aws_kms_key.getAwsKmsKeyPolicy", "apin_error", err)
		return nil, err
	}
	return keyPolicy, nil
}
