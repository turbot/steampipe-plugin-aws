package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/accessanalyzer"
	"github.com/aws/aws-sdk-go-v2/service/accessanalyzer/types"
	"github.com/aws/smithy-go"

	accessanalyzerv1 "github.com/aws/aws-sdk-go/service/accessanalyzer"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

type accessanalyzerFindingInfo = struct {
	Finding           types.Finding
	AccessAnalyzerArn string
}

func tableAwsAccessAnalyzerFinding(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_accessanalyzer_finding",
		Description: "AWS Access Analyzer Finding",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id", "access_analyzer_arn"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "ValidationException"}),
			},
			Hydrate: getAccessAnalyzerFinding,
			Tags:    map[string]string{"service": "access-analyzer", "action": "GetFinding"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAccessAnalyzers,
			Hydrate:       listAccessAnalyzersFindings,
			Tags:          map[string]string{"service": "access-analyzer", "action": "ListFindings"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException", "ValidationException"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "access_analyzer_arn",
					Require: plugin.Optional,
				},
				{
					Name:    "id",
					Require: plugin.Optional,
				},
				{
					Name:    "resource",
					Require: plugin.Optional,
				},
				{
					Name:    "status",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(accessanalyzerv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "access_analyzer_arn",
				Description: "The Amazon Resource Name (ARN) of the analyzer that generated the finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The ID of the finding.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Finding.Id"),
			},
			{
				Name:        "analyzed_at",
				Description: "The time at which the resource-based policy that generated the finding was analyzed.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Finding.AnalyzedAt"),
			},
			{
				Name:        "created_at",
				Description: "The time at which the finding was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Finding.CreatedAt"),
			},
			{
				Name:        "error",
				Description: "The error that resulted in an Error finding.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Finding.Error"),
			},
			{
				Name:        "is_public",
				Description: "Indicates whether the finding reports a resource that has a policy that allows public access.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Finding.IsPublic"),
			},
			{
				Name:        "resource",
				Description: "The resource that the external principal has access to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Finding.Resource"),
			},
			{
				Name:        "resource_owner_account",
				Description: "The Amazon Web Services account ID that owns the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Finding.ResourceOwnerAccount"),
			},
			{
				Name:        "resource_type",
				Description: "The type of the resource that the external principal has access to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Finding.ResourceType"),
			},
			{
				Name:        "status",
				Description: "The status of the finding.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Finding.Status"),
			},
			{
				Name:        "updated_at",
				Description: "The time at which the finding was most recently updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Finding.UpdatedAt"),
			},
			{
				Name:        "action",
				Description: "The action in the analyzed policy statement that an external principal has permission to use.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Finding.Action"),
			},
			{
				Name:        "sources",
				Description: "The sources of the finding, indicating how the access that generated the finding is granted. It is populated for Amazon S3 bucket findings.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Finding.Sources"),
			},
			{
				Name:        "principal",
				Description: "The external principal that has access to a resource within the zone of trust.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Finding.Principal"),
			},
			{
				Name:        "condition",
				Description: "The condition in the analyzed policy statement that resulted in a finding.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Finding.Condition"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Finding.Id"),
			},
		}),
	}
}

//// LIST FUNCTION

func listAccessAnalyzersFindings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var arn string
	if h.Item != nil {
		arn = *h.Item.(types.AnalyzerSummary).Arn
	}

	// Minimize API call with given Access analyzer ARN
	if arn != "" && d.EqualsQualString("access_analyzer_arn") != "" {
		if d.EqualsQualString("access_analyzer_arn") != arn {
			return nil, nil
		}
	}

	// Create session
	svc, err := AccessAnalyzerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_accessanalyzer_finding.listAccessAnalyzersFindings", "client_error", err)
		return nil, err
	}

	// The maximum number for MaxResults parameter is not defined by the API
	// We have set the MaxResults to 1000 based on our test
	maxItems := int32(1000)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = int32(limit)
		}
	}

	input := &accessanalyzer.ListFindingsInput{
		AnalyzerArn: &arn,
		MaxResults:  &maxItems,
		Filter:      map[string]types.Criterion{},
	}

	// set optional params
	setFilterCriteria(d, input)

	paginator := accessanalyzer.NewListFindingsPaginator(svc, input, func(o *accessanalyzer.ListFindingsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)

		if err != nil {
			var ae smithy.APIError
			if errors.As(err, &ae) {
				if ae.ErrorCode() == "ResourceNotFoundException" || ae.ErrorCode() == "ValidationException" {
					return nil, nil
				}
			}
			plugin.Logger(ctx).Error("aws_accessanalyzer_finding.listAccessAnalyzersFindings", "api_error", err)
			return nil, err
		}

		for _, finding := range output.Findings {
			f := types.Finding{
				AnalyzedAt:           finding.AnalyzedAt,
				Condition:            finding.Condition,
				CreatedAt:            finding.CreatedAt,
				Error:                finding.Error,
				Id:                   finding.Id,
				IsPublic:             finding.IsPublic,
				Principal:            finding.Principal,
				Resource:             finding.Resource,
				ResourceOwnerAccount: finding.ResourceOwnerAccount,
				ResourceType:         finding.ResourceType,
				Sources:              finding.Sources,
				Status:               finding.Status,
				UpdatedAt:            finding.UpdatedAt,
			}
			d.StreamListItem(ctx, accessanalyzerFindingInfo{f, arn})

			// Context may get canceled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAccessAnalyzerFinding(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQualString("id")
	arn := d.EqualsQualString("access_analyzer_arn")

	// check if id or arn is empty
	if id == "" || arn == "" {
		return nil, nil
	}

	// Create Session
	svc, err := AccessAnalyzerClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_accessanalyzer_finding.getAccessAnalyzerFinding", "client_error", err)
		return nil, err
	}

	// Build the params
	params := &accessanalyzer.GetFindingInput{
		AnalyzerArn: &arn,
		Id:          &id,
	}

	// Get call
	data, err := svc.GetFinding(ctx, params)

	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "ResourceNotFoundException" || ae.ErrorCode() == "ValidationException" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_accessanalyzer_finding.getAccessAnalyzerFinding", "api_error", err)
		return nil, err
	}

	return accessanalyzerFindingInfo{*data.Finding, arn}, nil
}

func setFilterCriteria(d *plugin.QueryData, input *accessanalyzer.ListFindingsInput) {
	params := []string{"id", "status", "resource"}
	for _, param := range params {
		if d.Quals[param] != nil {
			input.Filter[param] = types.Criterion{
				Eq: []string{d.EqualsQualString(param)},
			}
		}
	}
}
