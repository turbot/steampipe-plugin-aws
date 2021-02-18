package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsIamPolicySimulator(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_policy_simulator",
		Description: "AWS IAM Policy Simulator",
		List: &plugin.ListConfig{
			KeyColumns: plugin.AllColumns([]string{"principal_arn", "action", "resource_arn"}),
			Hydrate:    listIamPolicySimulation,
		},
		Columns: []*plugin.Column{
			// "Key" Columns
			{
				Name:        "principal_arn",
				Description: "The principal Amazon Resource Name (ARN) for this policy simulation",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "action",
				Description: "The action for this policy simulation",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "resource_arn",
				Type:        proto.ColumnType_STRING,
				Description: "The resource for this policy simulation",
				Transform:   transform.FromGo(),
			},
			{
				Name:        "decision",
				Type:        proto.ColumnType_STRING,
				Description: "The decision for this policy simulation",
				Transform:   transform.FromGo(),
			},
			{
				Name:        "decision_details",
				Type:        proto.ColumnType_JSON,
				Description: "The decision details for this policy simulation",
				Transform:   transform.FromGo(),
			},
			{
				Name:        "matched_statements",
				Type:        proto.ColumnType_JSON,
				Description: "The matched statements for this policy simulation",
				Transform:   transform.FromGo(),
			},
			{
				Name:        "missing_context_values",
				Type:        proto.ColumnType_JSON,
				Description: "The missing content values for this policy simulation",
				Transform:   transform.FromGo(),
			},
			{
				Name:        "resource_specific_results",
				Type:        proto.ColumnType_JSON,
				Description: "The resource specific results for this policy simulation",
				Transform:   transform.FromGo(),
			},
			{
				Name:        "organizations_decision_detail",
				Type:        proto.ColumnType_JSON,
				Description: "The organizations decision detail for this policy simulation",
				Transform:   transform.FromGo(),
			},
			{
				Name:        "permissions_boundary_decision_detail",
				Type:        proto.ColumnType_JSON,
				Description: "The permissions boundary decision detail for this policy simulation",
				Transform:   transform.FromGo(),
			},
		},
	}
}

type awsIamPolicySimulatorResult struct {
	Action                            string
	Decision                          *string
	DecisionDetails                   map[string]*string
	MatchedStatements                 []*iam.Statement
	MissingContextValues              []*string
	OrganizationsDecisionDetail       *iam.OrganizationsDecisionDetail
	PermissionsBoundaryDecisionDetail *iam.PermissionsBoundaryDecisionDetail
	PrincipalArn                      string
	ResourceArn                       string
	ResourceSpecificResults           []*iam.ResourceSpecificResult
	Result                            *iam.EvaluationResult
}

func listIamPolicySimulation(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listIamPolicySimulation")
	principalArn := d.KeyColumnQuals["principal_arn"].GetStringValue()
	action := d.KeyColumnQuals["action"].GetStringValue()
	resourceArn := d.KeyColumnQuals["resource_arn"].GetStringValue()

	// Create Session
	svc, err := IAMService(ctx, d)
	if err != nil {
		return nil, err
	}

	var params = &iam.SimulatePrincipalPolicyInput{
		PolicySourceArn: &principalArn,
		ActionNames:     []*string{&action},
		ResourceArns:    []*string{&resourceArn},
	}
	op, err := svc.SimulatePrincipalPolicy(params)
	if err != nil {
		return nil, err
	}

	evalResults := op.EvaluationResults
	resultForAction := evalResults[0]

	row := awsIamPolicySimulatorResult{
		Action:                            action,
		Decision:                          resultForAction.EvalDecision,
		DecisionDetails:                   resultForAction.EvalDecisionDetails,
		MatchedStatements:                 resultForAction.MatchedStatements,
		MissingContextValues:              resultForAction.MissingContextValues,
		OrganizationsDecisionDetail:       resultForAction.OrganizationsDecisionDetail,
		PermissionsBoundaryDecisionDetail: resultForAction.PermissionsBoundaryDecisionDetail,
		PrincipalArn:                      principalArn,
		ResourceArn:                       resourceArn,
		ResourceSpecificResults:           resultForAction.ResourceSpecificResults,
		Result:                            resultForAction,
	}

	d.StreamListItem(ctx, row)

	return nil, nil
}
