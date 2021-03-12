package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsGlacierVault(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_glacier_vault",
		Description: "AWS Glacier Vault",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("vault_name"),
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFoundException"}),
			Hydrate:           getGlacierVault,
		},
		List: &plugin.ListConfig{
			Hydrate: listGlacierVault,
		},
		GetMatrixItem: BuildRegionList,
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
				Name:        "tags_src",
				Description: "A list of tags associated with the cluster.",
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
				Transform:   transform.FromField("Tags"),
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
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// Create Session
	svc, err := GlacierService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)
	accountID := commonColumnData.AccountId

	// List call
	err = svc.ListVaultsPages(
		&glacier.ListVaultsInput{
			AccountId: aws.String(accountID),
		},
		func(page *glacier.ListVaultsOutput, isLast bool) bool {
			for _, vaults := range page.VaultList {
				d.StreamListItem(ctx, vaults)
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getGlacierVault(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	quals := d.KeyColumnQuals
	vaultName := quals["vault_name"].GetStringValue()

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)
	accountID := commonColumnData.AccountId

	// create service
	svc, err := GlacierService(ctx, d, region)
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

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	vaultName := h.Item.(*glacier.DescribeVaultOutput)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)
	accountID := commonColumnData.AccountId

	// Create session
	svc, err := GlacierService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build param
	param := &glacier.GetVaultAccessPolicyInput{
		VaultName: vaultName.VaultName,
		AccountId: aws.String(accountID),
	}

	vaultAccessPolicy, err := svc.GetVaultAccessPolicy(param)
	if err != nil {
		return nil, err
	}
	return vaultAccessPolicy, nil
}

func listTagsForGlacierVault(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listTagsForGlacierVault")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	vaultName := h.Item.(*glacier.DescribeVaultOutput)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}

	commonColumnData := commonData.(*awsCommonColumnData)
	accountID := commonColumnData.AccountId

	// Create session
	svc, err := GlacierService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build param
	param := &glacier.ListTagsForVaultInput{
		VaultName: vaultName.VaultName,
		AccountId: aws.String(accountID),
	}

	vaultTags, err := svc.ListTagsForVault(param)

	if err != nil {
		return nil, err
	}
	return vaultTags, nil
}
