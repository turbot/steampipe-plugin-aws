package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/backup"
	"github.com/aws/aws-sdk-go-v2/service/backup/types"

	backupv1 "github.com/aws/aws-sdk-go/service/backup"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsBackupLegalHold(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_backup_legal_hold",
		Description: "AWS Backup Legal Hold",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("legal_hold_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidParameterValueException"}),
			},
			Hydrate: getAwsBackupLegalHold,
			Tags:    map[string]string{"service": "backup", "action": "GetLegalHold"},
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsBackupLegalHolds,
			Tags:    map[string]string{"service": "backup", "action": "ListLegalHolds"},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(backupv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "legal_hold_id",
				Description: "ID of specific legal hold on one or more recovery points.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "This is an Amazon Resource Number (ARN) that uniquely identifies the legal hold.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LegalHoldArn"),
			},
			{
				Name:        "creation_date",
				Description: "This is the time in number format when legal hold was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "status",
				Description: "This is the status of the legal hold. Statuses can be ACTIVE, CREATING, CANCELED, and CANCELING.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cancellation_date",
				Description: "This is the time in number format when legal hold was cancelled.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "This is the description of a legal hold.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "retain_record_until",
				Description: "This is the date and time until which the legal hold record will be retained.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getAwsBackupLegalHold,
			},
			{
				Name:        "recovery_point_selection",
				Description: "This specifies criteria to assign a set of resources, such as resource types or backup vaults.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsBackupLegalHold,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LegalHoldArn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsBackupLegalHolds(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := BackupClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_legal_hold.listAwsBackupLegalHolds", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	input := &backup.ListLegalHoldsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	paginator := backup.NewListLegalHoldsPaginator(svc, input, func(o *backup.ListLegalHoldsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)
		
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_backup_legal_hold.listAwsBackupLegalHolds", "api_error", err)
			return nil, err
		}

		for _, item := range output.LegalHolds {
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsBackupLegalHold(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := BackupClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_legal_hold.getAwsBackupLegalHold", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	var id string
	if h.Item != nil {
		backupHold := h.Item.(types.LegalHold)
		id = *backupHold.LegalHoldId
	} else {
		id = d.EqualsQuals["legal_hold_id"].GetStringValue()
	}

	// check if id is empty
	if id == "" {
		return nil, nil
	}

	params := &backup.GetLegalHoldInput{
		LegalHoldId: aws.String(id),
	}

	op, err := svc.GetLegalHold(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_backup_legal_hold.getAwsBackupLegalHold", "api_error", err)
	}

	return op, nil
}
