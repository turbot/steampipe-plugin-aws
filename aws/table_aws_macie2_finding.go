package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/macie2"
	"github.com/aws/aws-sdk-go-v2/service/macie2/types"

	macie2v1 "github.com/aws/aws-sdk-go/service/macie2"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsMacie2Finding(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_macie2_finding",
		Description: "AWS Macie2 Finding",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ValidationException", "InvalidParameter"}),
			},
			Hydrate: getMacie2Finding,
			Tags:    map[string]string{"service": "macie2", "action": "GetFindings"},
		},
		List: &plugin.ListConfig{
			Hydrate: listMacie2Findings,
			Tags:    map[string]string{"service": "macie2", "action": "ListFindings"},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "type", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "severity", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "status", Require: plugin.Optional, Operators: []string{"=", "<>"}},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(macie2v1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The unique identifier for the finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of finding.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMacie2Finding,
			},
			{
				Name:        "severity",
				Description: "The severity level of the finding.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMacie2Finding,
			},
			{
				Name:        "status",
				Description: "The status of the finding.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMacie2Finding,
			},
			{
				Name:        "created_at",
				Description: "The date and time, in UTC and extended ISO 8601 format, when the finding was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getMacie2Finding,
			},
			{
				Name:        "updated_at",
				Description: "The date and time, in UTC and extended ISO 8601 format, when the finding was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getMacie2Finding,
				Transform:   transform.FromField("UpdatedAt").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "description",
				Description: "A description of the finding.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMacie2Finding,
			},
			{
				Name:        "category",
				Description: "The category of the finding.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMacie2Finding,
			},
			{
				Name:        "count",
				Description: "The total number of occurrences of the finding.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getMacie2Finding,
			},
			{
				Name:        "resources_affected",
				Description: "The resources that the finding applies to.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMacie2Finding,
			},
			{
				Name:        "sample",
				Description: "A sample of the data that triggered the finding.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMacie2Finding,
			},
			{
				Name:        "classification_details",
				Description: "The details of the classification that produced the finding.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMacie2Finding,
			},
			{
				Name:        "policy_details",
				Description: "The details of a policy finding. This value is null for a sensitive data finding.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMacie2Finding,
			},
			{
				Name:        "remediation",
				Description: "Information about the remediation steps for the finding.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMacie2Finding,
			},
			{
				Name:        "archived",
				Description: "Specifies whether the finding is archived.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getMacie2Finding,
			},
			{
				Name:        "source_region",
				Description: "The AWS Region where the finding was generated.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMacie2Finding,
			},
			{
				Name:        "source_account_id",
				Description: "The AWS account ID for the account that owns the finding.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMacie2Finding,
				Transform:   transform.FromField("AccountId"),
			},
			{
				Name:        "source_partition",
				Description: "The Amazon Web Services partition that Amazon Macie created the finding in.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMacie2Finding,
				Transform:   transform.FromField("Partition"),
			},
			{
				Name:        "source_region",
				Description: "The Amazon Web Services Region that Amazon Macie created the finding in.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMacie2Finding,
				Transform:   transform.FromField("Region"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMacie2Finding,
			},
		}),
	}
}

//// LIST FUNCTION

func listMacie2Findings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := Macie2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_macie2_finding.listMacie2Findings", "client_error", err)
		return nil, err
	}
	// Service is not supported in the region
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	maxItems := int32(50)
	input := &macie2.ListFindingsInput{
		MaxResults: aws.Int32(maxItems),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = int32(limit)
		}
	}

	filterCriteria := buildMacie2FindingsFilterCriteria(d.Quals)

	if len(filterCriteria.Criterion) > 0 {
		input.FindingCriteria = filterCriteria
	}

	// List call
	paginator := macie2.NewListFindingsPaginator(svc, input, func(o *macie2.ListFindingsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			// Throws "AccessDeniedException: Macie is not enabled." when AWS Macie is not enabled in a region
			// also the API throws AccessDeniedException if the request does not have proper permission
			// with the below check we will only handle "Macie is not enabled"
			if strings.Contains(err.Error(), "Macie is not enabled.") {
				return nil, nil
			}
			plugin.Logger(ctx).Error("aws_macie2_finding.listMacie2Findings", "api_error", err)
			return nil, err
		}

		for _, findingId := range output.FindingIds {
			d.StreamListItem(ctx, types.Finding{Id: &findingId})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getMacie2Finding(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var id string
	if h.Item != nil {
		id = *h.Item.(types.Finding).Id
	} else {
		id = d.EqualsQuals["id"].GetStringValue()
	}

	// empty check for finding id
	if id == "" {
		return nil, nil
	}

	// Create session
	svc, err := Macie2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_macie2_finding.getMacie2Finding", "client_error", err)
		return nil, err
	}
	// Service is not supported in the region
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build params
	params := &macie2.GetFindingsInput{
		FindingIds: []string{id},
	}

	// Get call
	op, err := svc.GetFindings(ctx, params)
	if err != nil {
		// Throws "AccessDeniedException: Macie is not enabled." when AWS Macie is not enabled in a region
		// also the API throws AccessDeniedException if the request does not have proper permission
		// with the below check we will only handle "Macie is not enabled"
		if strings.Contains(err.Error(), "Macie is not enabled.") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_macie2_finding.getMacie2Finding", "api_error", err)
		return nil, err
	}

	if len(op.Findings) > 0 {
		return op.Findings[0], nil
	}

	return nil, nil
}

//// UTILITY FUNCTION
//// Build macie2 list findings filter criteria

func buildMacie2FindingsFilterCriteria(quals plugin.KeyColumnQualMap) *types.FindingCriteria {
	filterCriteria := &types.FindingCriteria{
		Criterion: make(map[string]types.CriterionAdditionalProperties),
	}

	filterQuals := map[string]string{
		"finding_type": "type",
		"severity":     "severity",
		"status":       "status",
	}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			value := getQualsValueByColumn(quals, columnName, "string")

			filter := types.CriterionAdditionalProperties{
				Eq: []string{value.(string)},
			}

			if filterName == "type" {
				filterCriteria.Criterion["type"] = filter
			}
			if filterName == "severity" {
				filterCriteria.Criterion["severity"] = filter
			}
			if filterName == "status" {
				filterCriteria.Criterion["status"] = filter
			}
		}
	}

	return filterCriteria
}
