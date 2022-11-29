package aws

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glacier"
	"github.com/aws/aws-sdk-go-v2/service/glacier/types"
	"github.com/aws/smithy-go"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsGlacierVault(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_glacier_vault",
		Description: "AWS Glacier Vault",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("vault_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "InvalidParameter"}),
			},
			Hydrate: getGlacierVault,
		},
		List: &plugin.ListConfig{
			Hydrate: listGlacierVault,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "vault_name",
				Description: "The name of the vault.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vault_arn",
				Description: "The Amazon Resource Name (ARN) of the vault.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VaultARN"),
			},
			{
				Name:        "creation_date",
				Description: "The Universal Coordinated Time (UTC) date when the vault was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_inventory_date",
				Description: "The Universal Coordinated Time (UTC) date when Amazon S3 Glacier completed the last vault inventory.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "number_of_archives",
				Description: "The number of archives in the vault as of the last inventory date.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "size_in_bytes",
				Description: "Total size, in bytes, of the archives in the vault as of the last inventory date.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "policy",
				Description: "Contains the returned vault access policy as a JSON string.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGlacierVaultAccessPolicy,
				Transform:   transform.FromField("Policy.Policy"),
			},
			{
				Name:        "policy_std",
				Description: "Contains the policy in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGlacierVaultAccessPolicy,
				Transform:   transform.FromField("Policy.Policy").Transform(unescape).Transform(policyToCanonical),
			},
			{
				Name:        "vault_lock_policy",
				Description: "The vault lock policy.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGlacierVaultLockPolicy,
				Transform:   transform.FromField("Policy"),
			},
			{
				Name:        "vault_lock_policy_std",
				Description: "Contains the policy in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGlacierVaultLockPolicy,
				Transform:   transform.FromField("Policy").Transform(unescape).Transform(policyToCanonical),
			},
			{
				Name:        "vault_notification_config",
				Description: "Contains the notification configuration set on the vault.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGlacierVaultNotifications,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags associated with the vault.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForGlacierVault,
				Transform:   transform.FromField("Tags"),
			},

			// Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VaultName"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     listTagsForGlacierVault,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VaultARN").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listGlacierVault(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := GlacierClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glacier_vault.listGlacierVault", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glacier_vault.listGlacierVault", "api_error", err)
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)
	accountID := commonColumnData.AccountId
	maxLimit := int32(10)
	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(maxLimit) {
			if *limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = int32(*limit)
			}
		}
	}
	input := &glacier.ListVaultsInput{
		AccountId: aws.String(accountID),
		Limit:     aws.Int32(maxLimit),
	}
	paginator := glacier.NewListVaultsPaginator(svc, input, func(o *glacier.ListVaultsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})
	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_glacier_vault.listGlacierVault", "api_error", err)
			return nil, err
		}
		for _, vaults := range output.VaultList {
			d.StreamListItem(ctx, vaults)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getGlacierVault(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	vaultName := quals["vault_name"].GetStringValue()

	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glacier_vault.getGlacierVault", "api_error", err)
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)
	accountID := commonColumnData.AccountId

	// create service
	svc, err := GlacierClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glacier_vault.getGlacierVault", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	params := &glacier.DescribeVaultInput{
		VaultName: aws.String(vaultName),
		AccountId: aws.String(accountID),
	}

	op, err := svc.DescribeVault(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glacier_vault.getGlacierVault", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getGlacierVaultAccessPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := glacierVaultData(h.Item)
	accountID := strings.Split(data["Arn"], ":")[4]

	// Create session
	svc, err := GlacierClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glacier_vault.getGlacierVaultAccessPolicy", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build param
	param := &glacier.GetVaultAccessPolicyInput{
		VaultName: aws.String(data["Name"]),
		AccountId: aws.String(accountID),
	}

	vaultAccessPolicy, err := svc.GetVaultAccessPolicy(ctx, param)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glacier_vault.getGlacierVaultAccessPolicy", "api_error", err)
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "ResourceNotFoundException" {
				return nil, nil
			}
		}
	}
	return vaultAccessPolicy, nil
}

func getGlacierVaultLockPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := glacierVaultData(h.Item)
	accountID := strings.Split(data["Arn"], ":")[4]

	// Create session
	svc, err := GlacierClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glacier_vault.getGlacierVaultLockPolicy", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build param
	param := &glacier.GetVaultLockInput{
		VaultName: aws.String(data["Name"]),
		AccountId: aws.String(accountID),
	}

	vaultLock, err := svc.GetVaultLock(ctx, param)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glacier_vault.getGlacierVaultLockPolicy", "api_error", err)
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "ResourceNotFoundException" {
				return nil, nil
			}
		}
	}
	return vaultLock, nil
}

func getGlacierVaultNotifications(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := glacierVaultData(h.Item)
	accountID := strings.Split(data["Arn"], ":")[4]

	// Create session
	svc, err := GlacierClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glacier_vault.getGlacierVaultNotifications", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build param
	param := &glacier.GetVaultNotificationsInput{
		VaultName: aws.String(data["Name"]),
		AccountId: aws.String(accountID),
	}

	vaultNotifications, err := svc.GetVaultNotifications(ctx, param)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glacier_vault.getGlacierVaultNotifications", "api_error", err)
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "ResourceNotFoundException" {
				return nil, nil
			}
		}
	}
	return vaultNotifications, nil
}

func listTagsForGlacierVault(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := glacierVaultData(h.Item)
	accountID := strings.Split(data["Arn"], ":")[4]

	// Create session
	svc, err := GlacierClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glacier_vault.listTagsForGlacierVault", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Build param
	param := &glacier.ListTagsForVaultInput{
		VaultName: aws.String(data["Name"]),
		AccountId: aws.String(accountID),
	}

	vaultTags, err := svc.ListTagsForVault(ctx, param)
	if err != nil {
		plugin.Logger(ctx).Error("aws_glacier_vault.listTagsForGlacierVault", "api_error", err)
		return nil, err
	}
	return vaultTags, nil
}

func glacierVaultData(item interface{}) map[string]string {
	data := map[string]string{}
	switch item := item.(type) {
	case types.DescribeVaultOutput:
		data["Arn"] = *item.VaultARN
		data["Name"] = *item.VaultName
	case *glacier.DescribeVaultOutput:
		data["Arn"] = *item.VaultARN
		data["Name"] = *item.VaultName
	}
	return data
}
