package aws

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/securityhub"
	"github.com/aws/aws-sdk-go-v2/service/securityhub/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSecurityHubFinding(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_securityhub_finding",
		Description: "AWS Security Hub Finding",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getSecurityHubFinding,
			Tags:       map[string]string{"service": "securityhub", "action": "GetFindings"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidAccessException"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listSecurityHubFindings,
			Tags:    map[string]string{"service": "securityhub", "action": "GetFindings"},
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "company_name", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "compliance_status", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "confidence", Require: plugin.Optional, Operators: []string{"=", ">=", "<="}},
				{Name: "criticality", Require: plugin.Optional, Operators: []string{"=", ">=", "<="}},
				{Name: "generator_id", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "product_arn", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "product_name", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "record_state", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "title", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "verification_state", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "workflow_status", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "source_account_id", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "created_at", Require: plugin.Optional, Operators: []string{"=", ">=", ">", "<=", "<"}},
				{Name: "updated_at", Require: plugin.Optional, Operators: []string{"=", ">=", ">", "<=", "<"}},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidAccessException"}),
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_SECURITYHUB_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The security findings provider-specific identifier for a finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the finding.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},
			{
				Name:        "company_name",
				Description: "The name of the company for the product that generated the finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "confidence",
				Description: "A finding's confidence. Confidence is defined as the likelihood that a finding accurately identifies the behavior or issue that it was intended to identify.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "created_at",
				Description: "Indicates when the security-findings provider created the potential security issue that a finding captured.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "compliance_status",
				Description: "The result of a compliance standards check.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Compliance.Status"),
			},
			{
				Name:        "updated_at",
				Description: "Indicates when the security-findings provider last updated the finding record.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "criticality",
				Description: "The level of importance assigned to the resources associated with the finding.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "description",
				Description: "A finding's description.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "first_observed_at",
				Description: "Indicates when the security-findings provider first observed the potential security issue that a finding captured.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "generator_id",
				Description: "The identifier for the solution-specific component (a discrete unit of logic) that generated a finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_observed_at",
				Description: "Indicates when the security-findings provider most recently observed the potential security issue that a finding captured.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "product_arn",
				Description: "The ARN generated by Security Hub that uniquely identifies a product that generates findings.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "product_name",
				Description: "The name of the product that generated the finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "record_state",
				Description: "The record state of a finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "schema_version",
				Description: "The schema version that a finding is formatted for.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_url",
				Description: "A URL that links to a page about the current finding in the security-findings provider's solution.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "processed_at",
				Description: "An ISO8601-formatted timestamp that indicates when Security Hub received a finding and begins to process it.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "sample",
				Description: "Indicates whether the finding is a sample finding.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "verification_state",
				Description: "Indicates the veracity of a finding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "workflow_status",
				Description: "The workflow status of a finding. Possible values are NEW, NOTIFIED, SUPPRESSED, RESOLVED.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Workflow.Status"),
			},
			{
				Name:        "standards_control_arn",
				Description: "The ARN of the security standard control.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(extractStandardControlArn),
			},
			{
				Name:        "action",
				Description: "Provides details about an action that affects or that was taken on a resource.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "compliance",
				Description: "This data type is exclusive to findings that are generated as the result of a check run against a specific rule in a supported security standard, such as CIS Amazon Web Services Foundations.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "finding_provider_fields",
				Description: "In a BatchImportFindings request, finding providers use FindingProviderFields to provide and update their own values for confidence, criticality, related findings, severity, and types.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "malware",
				Description: "A list of malware related to a finding.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "network",
				Description: "The details of network-related information about a finding.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "network_path",
				Description: "Provides information about a network path that is relevant to a finding. Each entry under NetworkPath represents a component of that path.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "note",
				Description: "A user-defined note added to a finding.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "patch_summary",
				Description: "Provides an overview of the patch compliance status for an instance against a selected compliance standard.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "process",
				Description: "The details of process-related information about a finding.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "product_fields",
				Description: "A data type where security-findings providers can include additional solution-specific details that aren't part of the defined AwsSecurityFinding format.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "related_findings",
				Description: "A list of related findings.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "generator_details",
				Description: "Provides metadata for the Amazon CodeGuru detector associated with a finding.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "threats",
				Description: "Details about the threat detected in a security finding and the file paths that were affected by the threat.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "types",
				Description: "One or more finding types in the format of namespace/category/classifier that classify a finding. Valid namespace values are: Software and Configuration Checks | TTPs | Effects | Unusual Behaviors | Sensitive Data Identifications.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "remediation",
				Description: "A data type that describes the remediation options for a finding.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "resources",
				Description: "A set of resource data types that describe the resources that the finding refers to.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "severity",
				Description: "A finding's severity.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "threat_intel_indicators",
				Description: "Threat intelligence details related to a finding.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "user_defined_fields",
				Description: "A list of name/value string pairs associated with the finding.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "vulnerabilities",
				Description: "Provides a list of vulnerabilities associated with the findings.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "source_account_id",
				Description: "The account id where the affected resource lives.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AwsAccountId"),
			},
			/// Steampipe standard columns
			{
				Name:        "title",
				Description: "A finding's title.",
				Type:        proto.ColumnType_STRING,
			},
		}),
	}
}

//// LIST FUNCTION

func listSecurityHubFindings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := SecurityHubClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securityhub_finding.listSecurityHubFindings", "client_error", err)
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
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	input := &securityhub.GetFindingsInput{
		MaxResults: aws.Int32(maxLimit),
	}

	findingsFilter := buildListFindingsParam(d.Quals)
	if findingsFilter != nil {
		input.Filters = findingsFilter
	}

	// List call
	paginator := securityhub.NewGetFindingsPaginator(svc, input, func(o *securityhub.GetFindingsPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			// Handle error for accounts that are not subscribed to AWS Security Hub
			if strings.Contains(err.Error(), "not subscribed") {
				return nil, nil
			}
			plugin.Logger(ctx).Error("aws_securityhub_finding.listSecurityHubFindings", "api_error", err)
			return nil, err
		}

		for _, finding := range output.Findings {
			d.StreamListItem(ctx, finding)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSecurityHubFinding(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	id := d.EqualsQuals["id"].GetStringValue()

	// Empty check
	if id == "" {
		return nil, nil
	}

	// Create session
	svc, err := SecurityHubClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_securityhub_finding.getSecurityHubFinding", "client_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Build the params
	params := &securityhub.GetFindingsInput{
		Filters: &types.AwsSecurityFindingFilters{
			Id: []types.StringFilter{
				{
					Comparison: "EQUALS",
					Value:      aws.String(id),
				},
			},
		},
	}

	// Get call
	op, err := svc.GetFindings(ctx, params)
	if err != nil {
		// Handle error for unsupported or inactive regions
		if strings.Contains(err.Error(), "not subscribed") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("aws_securityhub_finding.getSecurityHubFinding", "api_error", err)
		return nil, err
	}
	if len(op.Findings) > 0 {
		return op.Findings[0], nil
	}
	return nil, nil
}

// Build param for findings list call
func buildListFindingsParam(quals plugin.KeyColumnQualMap) *types.AwsSecurityFindingFilters {
	securityFindingsFilter := &types.AwsSecurityFindingFilters{}
	strFilter := types.StringFilter{}
	dateFilter := types.DateFilter{}
	timeFormat := "2006-01-02T15:04:05Z"

	strColumns := []string{"company_name", "compliance_status", "generator_id", "product_arn", "product_name", "record_state", "title", "verification_state", "workflow_state", "workflow_status", "source_account_id"}

	timeColumns := []string{"created_at", "updated_at"}

	for _, t := range timeColumns {
		if quals[t] == nil {
			continue
		}
		for _, q := range quals[t].Quals {
			value := q.Value.GetTimestampValue().AsTime().Format(timeFormat)
			if value == "" {
				continue
			}

			switch q.Operator {
			case "=", ">=", ">":
				dateFilter.Start = &value
				dateFilter.End = aws.String(time.Now().Format(timeFormat))
			case "<", "<=":
				dateFilter.End = &value
				st, err := time.Parse(timeFormat, value)
				if err != nil {
					panic("failed to parsing provided value " + value + " for " + t)
				}
				if t == "updated_at" {
					// Default to past 90 days.
					// https://docs.aws.amazon.com/securityhub/latest/userguide/securityhub-findings.html
					dateFilter.Start = aws.String(st.AddDate(0, 0, -90).Format(timeFormat))
				} else {
					// For the query "select * from aws_securityhub_finding where created_at <= current_timestamp - interval '31d'", we are setting the end time based on the query parameter.
					// The query doesn't explicitly mention a start date for the creation time.
					// AWS retains Security Hub findings updated within the last 90 days, regardless of when they were created.
					// There is no strict limit for the creation start time, but we have set it to the date when AWS introduced SecurityHub. The API will return an error if we will not set the Start anf End time all together.
					findingIntroduceTime := "2018-11-27T00:00:00Z"
					t, err := time.Parse(timeFormat, findingIntroduceTime)
					if err != nil {
						panic("failed to parse the introduced securityhub findings time")
					}
					dateFilter.Start = aws.String(t.Format(timeFormat))
				}

			}
		}
		switch t {
		case "created_at":
			securityFindingsFilter.CreatedAt = []types.DateFilter{dateFilter}
		case "updated_at":
			securityFindingsFilter.UpdatedAt = []types.DateFilter{dateFilter}
		}
	}

	for _, s := range strColumns {
		if quals[s] == nil {
			continue
		}
		for _, q := range quals[s].Quals {
			value := q.Value.GetStringValue()
			if value == "" {
				continue
			}

			switch q.Operator {
			case "<>":
				strFilter.Comparison = "NOT_EQUALS"
			case "=":
				strFilter.Comparison = "EQUALS"
			}

			switch s {
			case "company_name":
				strFilter.Value = aws.String(value)
				securityFindingsFilter.CompanyName = append(securityFindingsFilter.CompanyName, strFilter)
			case "generator_id":
				strFilter.Value = aws.String(value)
				securityFindingsFilter.GeneratorId = append(securityFindingsFilter.GeneratorId, strFilter)
			case "compliance_status":
				strFilter.Value = aws.String(value)
				securityFindingsFilter.ComplianceStatus = append(securityFindingsFilter.ComplianceStatus, strFilter)
			case "product_arn":
				strFilter.Value = aws.String(value)
				securityFindingsFilter.ProductArn = append(securityFindingsFilter.ProductArn, strFilter)
			case "product_name":
				strFilter.Value = aws.String(value)
				securityFindingsFilter.ProductName = append(securityFindingsFilter.ProductName, strFilter)
			case "record_state":
				strFilter.Value = aws.String(value)
				securityFindingsFilter.RecordState = append(securityFindingsFilter.RecordState, strFilter)
			case "title":
				strFilter.Value = aws.String(value)
				securityFindingsFilter.Title = append(securityFindingsFilter.Title, strFilter)
			case "verification_state":
				strFilter.Value = aws.String(value)
				securityFindingsFilter.VerificationState = append(securityFindingsFilter.VerificationState, strFilter)
			case "workflow_state":
				strFilter.Value = aws.String(value)
				securityFindingsFilter.WorkflowState = append(securityFindingsFilter.WorkflowState, strFilter)
			case "workflow_status":
				strFilter.Value = aws.String(value)
				securityFindingsFilter.WorkflowStatus = append(securityFindingsFilter.WorkflowStatus, strFilter)
			case "source_account_id":
				strFilter.Value = aws.String(value)
				securityFindingsFilter.AwsAccountId = append(securityFindingsFilter.AwsAccountId, strFilter)
			}

		}
	}

	return securityFindingsFilter
}

//// TRANSFORM FUNCTIONS

func extractStandardControlArn(_ context.Context, d *transform.TransformData) (interface{}, error) {
	findingArn := d.HydrateItem.(types.AwsSecurityFinding).Id

	if strings.Contains(*findingArn, "arn:aws:securityhub") {
		standardControlArn := strings.Replace(strings.Split(*findingArn, "/finding")[0], "subscription", "control", 1)
		return standardControlArn, nil
	}
	return nil, nil
}
