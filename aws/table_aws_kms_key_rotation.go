package aws

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/aws/smithy-go"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsKmsKeyRotation(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_kms_key_rotation",
		Description: "AWS KMS Key Rotation",
		List: &plugin.ListConfig{
			ParentHydrate: listKmsKeys,
			Hydrate:       listKmsKeyRotations,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NotFoundException"}),
			},
			Tags: map[string]string{"service": "kms", "action": "ListKeyRotations"},
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "key_id",
					Require: plugin.Optional,
				},
				{
					Name:    "key_arn",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_KMS_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "key_id",
				Description: "Unique identifier of the key.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "key_arn",
				Description: "ARN of the key.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "rotation_type",
				Description: "Identifies whether the key material rotation was a scheduled automatic rotation or an on-demand rotation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "rotation_date",
				Description: "Date and time that the key material rotation completed.",
				Type:        proto.ColumnType_TIMESTAMP,
			},

			/// Standard columns for all tables
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("KeyId"),
			},
		}),
	}
}

type RotationInfo struct {
	KeyId        *string
	KeyArn       *string
	RotationDate *time.Time
	RotationType types.RotationType
}

//// LIST FUNCTION

func listKmsKeyRotations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var key types.KeyListEntry
	if h.Item != nil {
		key = h.Item.(types.KeyListEntry)
	}

	if d.EqualsQualString("key_id") != "" || d.EqualsQualString("key_arn") != "" {
		if d.EqualsQualString("key_id") != "" && d.EqualsQualString("key_id") != *key.KeyId {
			return nil, nil
		}
		if d.EqualsQualString("key_arn") != "" && d.EqualsQualString("key_arn") != *key.KeyArn {
			return nil, nil
		}
	}

	// Create Session
	svc, err := KMSClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_kms_key_rotation.listKmsKeyRotations", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	maxItems := int32(1000)
	input := &kms.ListKeyRotationsInput{
		KeyId: key.KeyArn,
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
	paginator := kms.NewListKeyRotationsPaginator(svc, input, func(o *kms.ListKeyRotationsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			// In the case of parent hydrate the ignore config seems to not work for the child table. So we need to handle it manually.
			// Steampipe SDK issue ref: https://github.com/turbot/steampipe-plugin-sdk/issues/544
			ignoreCodes := GetConfig(d.Connection).IgnoreErrorCodes
			var ae smithy.APIError
			if errors.As(err, &ae) {
				if helpers.StringSliceContains(ignoreCodes, ae.ErrorCode()) {
					return nil, nil
				}
			}
			plugin.Logger(ctx).Error("aws_kms_key_rotation.listKmsKeyRotations", "api_error", err)
			return nil, err
		}

		for _, rotation := range output.Rotations {
			d.StreamListItem(ctx, &RotationInfo{
				KeyId:        aws.String(strings.Split(*rotation.KeyId, "/")[1]),
				KeyArn:       rotation.KeyId,
				RotationDate: rotation.RotationDate,
				RotationType: rotation.RotationType,
			})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
