package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wellarchitected"
	"github.com/aws/aws-sdk-go-v2/service/wellarchitected/types"
	"github.com/aws/smithy-go"

	wellarchitectedv1 "github.com/aws/aws-sdk-go/service/wellarchitected"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsWellArchitectedLensShare(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wellarchitected_lens_share",
		Description: "AWS Well-Architected Lens Share",
		List: &plugin.ListConfig{
			ParentHydrate: listWellArchitectedLenses,
			Hydrate:       listWellArchitectedLensShares,
			Tags:          map[string]string{"service": "wellarchitected", "action": "ListLensShares"},
			// TODO: Uncomment and remove extra check in
			// listWellArchitectedLensShares function once this works again
			// IgnoreConfig: &plugin.IgnoreConfig{
			// 	ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			// },
			KeyColumns: []*plugin.KeyColumn{
				{Name: "lens_alias", Require: plugin.Optional},
				{Name: "shared_with", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(wellarchitectedv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "lens_alias",
				Description: "The alias of the lens.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lens_name",
				Description: "The full name of the lens.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lens_arn",
				Description: "The ARN of the lens.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "share_id",
				Description: "The ID associated with the workload share.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "shared_with",
				Description: "The Amazon Web Services account ID, IAM role, organization ID, or organizational unit (OU) ID with which the workload is shared.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of a workload share.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status_message",
				Description: "Optional message to compliment the Status field.",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ShareId"),
			},
		}),
	}
}

type LensShareInfo struct {
	LensAlias     *string
	LensName      *string
	LensArn       *string
	ShareId       *string
	SharedWith    *string
	Status        types.ShareStatus
	StatusMessage *string
}

//// LIST FUNCTION

func listWellArchitectedLensShares(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	lens := h.Item.(types.LensSummary)
	lensAlias := d.EqualsQualString("lens_alias")

	// Reduce the number of API calls if the 'lens_alias' has been provided in the query parameter
	if lensAlias != "" {
		if lensAlias != *lens.LensAlias || lensAlias != *lens.LensArn {
			return nil, nil
		}
	}

	// Create session
	svc, err := WellArchitectedClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_wellarchitected_lens_share.listWellArchitectedLensShares", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// As per the doc(https://docs.aws.amazon.com/wellarchitected/latest/APIReference/API_ListLensShares.html) the 'MaxResults' value can be between 1-50.
	// But the API throws ValidationException error(Error: operation error WellArchitected: ListLensShares, https response error StatusCode: 400, RequestID: 90b963e0-4060-41a3-8081-2e3da2660f56, ValidationException: [Validation] {"reason":"FIELD_VALIDATION_FAILED","message":"[Validation] MaxResults must be between 1 and 10.","fields":"[{\"Name\":\"MaxResults\",\"Message\":\"MaxResults must be between 1 and 10.\"}]"})

	// Limiting the results
	maxLimit := int32(10)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &wellarchitected.ListLensSharesInput{
		MaxResults: aws.Int32(maxLimit),
		LensAlias:  lens.LensArn,
	}

	if d.EqualsQualString("shared_with") != "" {
		input.SharedWithPrefix = aws.String(d.EqualsQualString("shared_with"))
	}

	if d.EqualsQualString("status") != "" {
		input.Status = types.ShareStatus(d.EqualsQualString("status"))
	}

	paginator := wellarchitected.NewListLensSharesPaginator(svc, input, func(o *wellarchitected.ListLensSharesPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			var ae smithy.APIError
			if errors.As(err, &ae) {
				if ae.ErrorCode() == "ResourceNotFoundException" {
					return nil, nil
				}
			}
			plugin.Logger(ctx).Error("aws_wellarchitected_lens_share.listWellArchitectedLensShares", "api_error", err)
			return nil, err
		}

		for _, item := range output.LensShareSummaries {

			if lens.LensAlias == nil || *lens.LensAlias == "" {
				lens.LensAlias = lens.LensArn
			}
			d.StreamListItem(ctx, LensShareInfo{
				LensAlias:     lens.LensAlias,
				LensName:      lens.LensName,
				LensArn:       lens.LensArn,
				ShareId:       item.ShareId,
				SharedWith:    item.SharedWith,
				Status:        item.Status,
				StatusMessage: item.StatusMessage,
			})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, nil
}
