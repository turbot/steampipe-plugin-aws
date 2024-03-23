package aws

import (
	"context"
	"errors"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
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
					Name:    "action",
					Require: plugin.Optional,
				},
				{
					Name:    "analyzed_at",
					Require: plugin.Optional,
				},
				{
					Name:    "condition",
					Require: plugin.Optional,
				},
				{
					Name:    "created_at",
					Require: plugin.Optional,
				},
				{
					Name:    "error",
					Require: plugin.Optional,
				},
				{
					Name:    "is_public",
					Require: plugin.Optional,
				},
				{
					Name:    "principal",
					Require: plugin.Optional,
				},
				{
					Name:    "resource",
					Require: plugin.Optional,
				},
				{
					Name:    "resource_owner_account",
					Require: plugin.Optional,
				},
				{
					Name:    "resource_type",
					Require: plugin.Optional,
				},
				{
					Name:    "sources",
					Require: plugin.Optional,
				},
				{
					Name:    "status",
					Require: plugin.Optional,
				},
				{
					Name:    "updated_at",
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
				Name:        "action",
				Description: "The action in the analyzed policy statement that an external principal has permission to use.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Finding.Action"),
			},
			{
				Name:        "analyzed_at",
				Description: "The time at which the resource-based policy that generated the finding was analyzed.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Finding.AnalyzedAt"),
			},
			{
				Name:        "condition",
				Description: "The condition in the analyzed policy statement that resulted in a finding.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Finding.Condition"),
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
				Name:        "principal",
				Description: "The external principal that has access to a resource within the zone of trust.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Finding.Principal"),
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
				Name:        "sources",
				Description: "The sources of the finding, indicating how the access that generated the finding is granted. It is populated for Amazon S3 bucket findings.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Finding.Sources"),
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

	if arn == "" {
		return nil, nil
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

	filters := buildListFindingsFilters(d.Quals)

	input := &accessanalyzer.ListFindingsInput{
		AnalyzerArn: &arn,
		Filter:      filters,
		MaxResults:  &maxItems,
	}

	paginator := accessanalyzer.NewListFindingsPaginator(svc, input, func(o *accessanalyzer.ListFindingsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_accessanalyzer_finding.listAccessAnalyzersFindings", "api_error", err)
			return nil, err
		}

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

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAccessAnalyzerFinding(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQuals["id"].GetStringValue()
	arn := d.EqualsQuals["access_analyzer_arn"].GetStringValue()

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
		plugin.Logger(ctx).Debug("aws_accessanalyzer_finding.getAccessAnalyzerFinding", "api_error", err)
		return nil, err
	}

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

// Build the list call input filter
//// UTILITY FUNCTION

// Build Access Analyzer list call input filter
func buildListFindingsFilters(quals plugin.KeyColumnQualMap) map[string]types.Criterion {
	filters := make(map[string]types.Criterion)

	for qualName, qual := range quals {
			criterion := types.Criterion{}

			switch qualName {
			case "contains", "eq", "neq":
					var stringPointers []*string
					for _, q := range qual.Quals {
						// Convert each string to a *string and append to the slice
						stringVal := q.Value.GetStringValue() // Get the string value
						stringPointers = append(stringPointers, &stringVal) // Append the pointer to the slice
					}

					// Convert the slice of string pointers to a slice of strings
					var stringSlice []string
					for _, sp := range stringPointers {
						stringSlice = append(stringSlice, *sp)
					}

					// Directly assign the slice of strings to the appropriate field
					if qualName == "contains" {
						criterion.Contains = stringSlice
					} else if qualName == "eq" {
						criterion.Eq = stringSlice
					} else if qualName == "neq" {
						criterion.Neq = stringSlice
					}

			case "exists":
					if len(qual.Quals) > 0 {
							existsValue, err := strconv.ParseBool(qual.Quals[0].Value.GetStringValue())
							if err == nil {
									criterion.Exists = aws.Bool(existsValue) // Directly use aws.Bool for a boolean pointer
							}
					}
			}

			// Ensure criterion is not empty before adding to filters
			if len(criterion.Contains) > 0 || len(criterion.Eq) > 0 || len(criterion.Neq) > 0 || criterion.Exists != nil {
					filters[qualName] = criterion
			}
	}

	return filters
}
