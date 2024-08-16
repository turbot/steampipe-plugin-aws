package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/qldb"

	qldbv1 "github.com/aws/aws-sdk-go/service/qldb"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsQLDBLedger(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_qldb_ledger",
		Description: "AWS QLDB Ledger",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getAwsQldbLedger,
			Tags:    map[string]string{"service": "qldb", "action": "DescribeLedger"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsQldbLedgers,
			Tags:    map[string]string{"service": "qldb", "action": "ListLedgers"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getAwsQldbLedger,
				Tags: map[string]string{"service": "qldb", "action": "DescribeLedger"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(qldbv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the ledger.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the ledger.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsQldbLedger,
			},
			{
				Name:        "creation_time",
				Description: "The date and time, in epoch time format, when the ledger was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "state",
				Description: "The current status of the ledger.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "deletion_protection",
				Description: "Specifies whether the ledger is protected from being deleted by any user.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getAwsQldbLedger,
			},
			{
				Name:        "permissions_mode",
				Description: "The permissions mode of the ledger.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsQldbLedger,
			},
			{
				Name:        "encryption_status",
				Description: "The current state of encryption at rest for the ledger.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsQldbLedger,
				Transform:   transform.FromField("EncryptionDescription.EncryptionStatus"),
			},
			{
				Name:        "kms_key_arn",
				Description: "The Amazon Resource Name (ARN) of the customer managed KMS key that the ledger uses for encryption at rest.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsQldbLedger,
				Transform:   transform.FromField("EncryptionDescription.KmsKeyArn"),
			},
			{
				Name:        "inaccessible_kms_key_date_time",
				Description: "The date and time, in epoch time format, when the KMS key first became inaccessible, in the case of an error.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getAwsQldbLedger,
				Transform:   transform.FromField("EncryptionDescription.InaccessibleKmsKeyDateTime"),
			},

			// Steampipe Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsQldbLedger,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsQldbLedgers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get client
	svc, err := QLDBClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_qldb_ledger.listAwsQldbLedgers", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &qldb.ListLedgersInput{
		MaxResults: &maxLimit,
	}

	paginator := qldb.NewListLedgersPaginator(svc, input, func(o *qldb.ListLedgersPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_qldb_ledger.listAwsQldbLedgers", "api_error", err)
			return nil, err
		}

		for _, item := range output.Ledgers {
			d.StreamListItem(ctx, &qldb.DescribeLedgerOutput{
				Name:             item.Name,
				State:            item.State,
				CreationDateTime: item.CreationDateTime,
			})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsQldbLedger(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := QLDBClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_qldb_ledger.getAwsQldbLedger", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	name := d.EqualsQualString("name")
	if h.Item != nil {
		name = *h.Item.(*qldb.DescribeLedgerOutput).Name
	}

	// Empty Check
	if name == "" {
		return nil, nil
	}

	// Build the params
	params := &qldb.DescribeLedgerInput{
		Name: &name,
	}

	// Get call
	data, err := svc.DescribeLedger(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_qldb_ledger.getAwsQldbLedger", "api_error", err)
		return nil, err
	}

	return data, nil
}
