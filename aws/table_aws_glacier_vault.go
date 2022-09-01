package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableAwsGlacierVault(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_glacier_vault",
		Description: "AWS Glacier Vault",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("vault_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFoundException", "InvalidParameter"}),
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
	svc, err := GlacierService(ctx, d)
	if err != nil {
		return nil, err
	}

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)
	accountID := commonColumnData.AccountId
	maxLimit := "10"

	input := &glacier.ListVaultsInput{
		AccountId: aws.String(accountID),
		Limit:     &maxLimit,
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < 10 {
			if *limit < 1 {
				input.Limit = aws.String("1")
			} else {
				input.Limit = aws.String(fmt.Sprint(*limit))
			}
		}
	}

	// List call
	err = svc.ListVaultsPages(
		input,
		func(page *glacier.ListVaultsOutput, isLast bool) bool {
			for _, vaults := range page.VaultList {
				d.StreamListItem(ctx, vaults)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getGlacierVault(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	vaultName := quals["vault_name"].GetStringValue()

	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)
	accountID := commonColumnData.AccountId

	// create service
	svc, err := GlacierService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &glacier.DescribeVaultInput{
		VaultName: aws.String(vaultName),
		AccountId: aws.String(accountID),
	}

	op, err := svc.DescribeVault(params)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func getGlacierVaultAccessPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getGlacierVaultAccessPolicy")

	data := h.Item.(*glacier.DescribeVaultOutput)
	accountID := strings.Split(*data.VaultARN, ":")[4]

	// Create session
	svc, err := GlacierService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build param
	param := &glacier.GetVaultAccessPolicyInput{
		VaultName: data.VaultName,
		AccountId: aws.String(accountID),
	}

	vaultAccessPolicy, err := svc.GetVaultAccessPolicy(param)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "ResourceNotFoundException" {
				return nil, nil
			}
			return nil, err
		}
	}
	return vaultAccessPolicy, nil
}

func getGlacierVaultLockPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getGlacierVaultLockPolicy")

	data := h.Item.(*glacier.DescribeVaultOutput)
	accountID := strings.Split(*data.VaultARN, ":")[4]

	// Create session
	svc, err := GlacierService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build param
	param := &glacier.GetVaultLockInput{
		VaultName: data.VaultName,
		AccountId: aws.String(accountID),
	}

	vaultLock, err := svc.GetVaultLock(param)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "ResourceNotFoundException" {
				return nil, nil
			}
			return nil, err
		}
	}
	return vaultLock, nil
}

func getGlacierVaultNotifications(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	data := h.Item.(*glacier.DescribeVaultOutput)
	accountID := strings.Split(*data.VaultARN, ":")[4]

	// Create session
	svc, err := GlacierService(ctx, d)
	if err != nil {
		logger.Error("aws_glacier_vault.getGlacierVaultNotifications", "service_creation_error", err)
		return nil, err
	}

	// Build param
	param := &glacier.GetVaultNotificationsInput{
		VaultName: data.VaultName,
		AccountId: aws.String(accountID),
	}

	vaultNotifications, err := svc.GetVaultNotifications(param)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "ResourceNotFoundException" {
				return nil, nil
			}
			logger.Error("aws_glacier_vault.getGlacierVaultNotifications", "api_error", err)
			return nil, err
		}
	}
	return vaultNotifications, nil
}

func listTagsForGlacierVault(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listTagsForGlacierVault")

	data := h.Item.(*glacier.DescribeVaultOutput)
	accountID := strings.Split(*data.VaultARN, ":")[4]

	// Create session
	svc, err := GlacierService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build param
	param := &glacier.ListTagsForVaultInput{
		VaultName: data.VaultName,
		AccountId: aws.String(accountID),
	}

	vaultTags, err := svc.ListTagsForVault(param)
	if err != nil {
		return nil, err
	}
	return vaultTags, nil
}
