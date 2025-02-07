package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"

	ssmEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsSSMDocumentPermission(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssm_document_permission",
		Description: "AWS SSM Document Permission",
		List: &plugin.ListConfig{
			Hydrate: listAwsSSMDocumentPermissions,
			Tags:    map[string]string{"service": "ssm", "action": "DescribeDocumentPermission"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "document_name", Require: plugin.Required},
			},
			// Getting InvalidDocument error if document is not avaiable in any specific region
			// Error: aws: operation error SSM: DescribeDocumentPermission, https response error StatusCode: 400, RequestID: fece9f20-9b41-40d9-abd3-933c1c1e4345, InvalidDocument: Document with name ssm_doc_test_delete does not exist. (SQLSTATE HV000)
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidDocument"}),
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ssmEndpoint.AWS_SSM_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "document_name",
				Description: "The name of the Systems Manager document.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "shared_account_id",
				Description: "The Amazon Web Services account ID where the current document is shared.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountSharingInfo.AccountId"),
			},
			{
				Name:        "shared_document_version",
				Description: "The version of the current document shared with the account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountSharingInfo.SharedDocumentVersion"),
			},
			{
				Name:        "account_ids",
				Description: "The account IDs that have permission to use this document. The ID can be either an AWS account or All.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountSharingInfo.SharedDocumentVersion"),
			},
		}),
	}
}

type PermissionInfo struct {
	DocumentName       string
	AccountIds         []string
	AccountSharingInfo types.AccountSharingInfo
}

//// LIST FUNCTION

func listAwsSSMDocumentPermissions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	documentName := d.EqualsQualString("document_name")

	// Empty check
	if documentName == "" {
		return nil, nil
	}

	// Create session
	svc, err := SSMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ssm_document_permission.listAwsSSMDocumentPermissions", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	maxItems := int32(200)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = int32(limit)
		}
	}

	input := &ssm.DescribeDocumentPermissionInput{
		MaxResults:     &maxItems,
		Name:           &documentName,
		PermissionType: types.DocumentPermissionType("Share"),
	}

	// API doesn't support aws-sdk-go-v2 paginator as of date
	pagesLeft := true
	for pagesLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		response, err := svc.DescribeDocumentPermission(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("aws_ssm_document_permission.listAwsSSMDocumentPermissions", "api_error", err)
			return nil, err
		}
		for _, item := range response.AccountSharingInfoList {
			d.StreamListItem(ctx, &PermissionInfo{
				DocumentName:       documentName,
				AccountIds:         response.AccountIds,
				AccountSharingInfo: item,
			})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if response.NextToken != nil {
			pagesLeft = true
			input.NextToken = response.NextToken
		} else {
			pagesLeft = false
		}
	}

	return nil, nil
}
