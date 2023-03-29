package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"

	kmsv1 "github.com/aws/aws-sdk-go/service/kms"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsKmsAlias(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_kms_alias",
		Description: "AWS KMS Alias",
		List: &plugin.ListConfig{
			ParentHydrate: listKmsKeys,
			Hydrate:       listKmsAliases,
		},
		GetMatrixItemFunc: SupportedRegionMatrix(kmsv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "alias_name",
				Description: "String that contains the alias. This value begins with alias/.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "String that contains the key ARN.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AliasArn"),
			},
			{
				Name:        "target_key_id",
				Description: "String that contains the key identifier of the KMS key associated with the alias.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_date",
				Description: "Date and time that the alias was most recently created in the account and Region.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_updated_date",
				Description: "Date and time that the alias was most recently associated with a KMS key in the account and Region.",
				Type:        proto.ColumnType_TIMESTAMP,
			},

			/// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AliasName"),
			},

			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AliasArn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listKmsAliases(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	keyInfo := h.Item.(types.KeyListEntry)
	keyId := keyInfo.KeyId

	// Create Client
	svc, err := KMSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kms_key_alias.listKmsAliases", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	maxItems := int32(100)
	input := &kms.ListAliasesInput{
		KeyId: keyId,
	}

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
	paginator := kms.NewListAliasesPaginator(svc, input, func(o *kms.ListAliasesPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_kms_key_alias.listKmsAliases", "api_error", err)
			return nil, err
		}

		for _, alias := range output.Aliases {
			d.StreamListItem(ctx, alias)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
