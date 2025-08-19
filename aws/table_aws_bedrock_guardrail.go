package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/bedrock"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// unified row used for both List and Get paths
type bedrockGuardrailRow struct {
	Arn         string    `json:"Arn"`
	GuardrailId string    `json:"GuardrailId"`
	Name        string    `json:"Name"`
	Description string    `json:"Description"`
	Status      string    `json:"Status"`
	Version     string    `json:"Version"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
}

func tableAwsBedrockGuardrail(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_bedrock_guardrail",
		Description: "Amazon Bedrock Guardrail.",
		List: &plugin.ListConfig{
			Hydrate: listBedrockGuardrails,
			Tags:    map[string]string{"service": "bedrock", "action": "ListGuardrails"},
		},
		Get: &plugin.GetConfig{
			// allow lookup by ID or ARN (both map to GuardrailIdentifier)
			KeyColumns: plugin.AnyColumn([]string{"guardrail_id", "arn"}),
			Hydrate:    getBedrockGuardrail,
			Tags:       map[string]string{"service": "bedrock", "action": "GetGuardrail"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_BEDROCK_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			// identifiers
			{Name: "arn", Type: proto.ColumnType_STRING, Description: "ARN of the guardrail.", Transform: transform.FromField("Arn")},
			{Name: "guardrail_id", Type: proto.ColumnType_STRING, Description: "ID of the guardrail.", Transform: transform.FromField("GuardrailId")},

			// metadata
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the guardrail.", Transform: transform.FromField("Name")},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of the guardrail.", Transform: transform.FromField("Description")},

			// status / version
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the guardrail.", Transform: transform.FromField("Status")},
			{Name: "version", Type: proto.ColumnType_STRING, Description: "Version (DRAFT or a number).", Transform: transform.FromField("Version")},

			// timestamps
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt")},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("UpdatedAt")},

			// steampipe standard
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name")},
			{Name: "akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("Arn").Transform(transform.EnsureStringArray)},
		}),
	}
}

// LIST: map GuardrailSummary -> bedrockGuardrailRow (ensures Arn/Id are set)
func listBedrockGuardrails(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	svc, err := BedrockClient(ctx, d)
	if svc == nil {
		return nil, nil
	}
	if err != nil {
		plugin.Logger(ctx).Error("aws_bedrock_guardrail.listBedrockGuardrails", "connection_error", err)
		return nil, err
	if err != nil {
		plugin.Logger(ctx).Error("aws_bedrock_guardrail.listBedrockGuardrails", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, err
	}

	p := bedrock.NewListGuardrailsPaginator(svc, &bedrock.ListGuardrailsInput{})
	for p.HasMorePages() {
		out, err := p.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_bedrock_guardrail.listBedrockGuardrails", "api_error", err)
			return nil, err
		}
		for _, s := range out.Guardrails {
			row := bedrockGuardrailRow{
				Arn:         str(s.Arn),
				GuardrailId: str(s.Id),
				Name:        str(s.Name),
				Description: str(s.Description),
				Status:      string(s.Status),
				Version:     str(s.Version),
				CreatedAt:   t(s.CreatedAt),
				UpdatedAt:   t(s.UpdatedAt),
			}
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, nil
}

// GET: map GetGuardrailOutput -> bedrockGuardrailRow (ensures Arn/Id are set)
func getBedrockGuardrail(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQualString("guardrail_id")
	if id == "" {
		id = d.EqualsQualString("arn")
	}
	if id == "" {
		return nil, nil
	}

	svc, err := BedrockClient(ctx, d)
	if svc == nil {
		return nil, nil
	}
	if err != nil {
		plugin.Logger(ctx).Error("aws_bedrock_guardrail.getBedrockGuardrail", "connection_error", err)
		return nil, err
	}

	out, err := svc.GetGuardrail(ctx, &bedrock.GetGuardrailInput{
		GuardrailIdentifier: &id, // accepts ID or ARN
	})
	if err != nil {
		plugin.Logger(ctx).Error("aws_bedrock_guardrail.getBedrockGuardrail", "api_error", err)
		return nil, err
	}

	row := bedrockGuardrailRow{
		Arn:         str(out.GuardrailArn),
		GuardrailId: str(out.GuardrailId),
		Name:        str(out.Name),
		Description: str(out.Description),
		Status:      string(out.Status),
		Version:     str(out.Version),
		CreatedAt:   t(out.CreatedAt),
		UpdatedAt:   t(out.UpdatedAt),
	}
	return row, nil
}

// small ptr helpers (avoid extra deps)
func str(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
func t(p *time.Time) time.Time {
	if p == nil {
		return time.Time{}
	}
	return *p
}


