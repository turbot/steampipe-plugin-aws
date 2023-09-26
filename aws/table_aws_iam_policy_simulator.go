package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsIamPolicySimulator(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_policy_simulator",
		Description: "AWS IAM Policy Simulator",
		List: &plugin.ListConfig{
			KeyColumns: plugin.AllColumns([]string{"principal_arn", "action", "resource_arn"}),
			Hydrate:    listIamPolicySimulation,
			Tags:       map[string]string{"service": "iam", "action": "SimulatePrincipalPolicy"},
		},
		Columns: []*plugin.Column{
			// "Key" Columns
			{
				Name:        "principal_arn",
				Description: "The principal Amazon Resource Name (ARN) for this policy simulation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "action",
				Description: "The action for this policy simulation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "resource_arn",
				Type:        proto.ColumnType_STRING,
				Description: "The resource for this policy simulation.",
				Transform:   transform.FromGo(),
			},
			{
				Name:        "decision",
				Type:        proto.ColumnType_STRING,
				Description: "The decision for this policy simulation.",
				Transform:   transform.FromGo(),
			},
			{
				Name:        "decision_details",
				Type:        proto.ColumnType_JSON,
				Description: "The decision details for this policy simulation.",
				Transform:   transform.FromGo(),
			},
			{
				Name:        "matched_statements",
				Type:        proto.ColumnType_JSON,
				Description: "The matched statements for this policy simulation.",
				Transform:   transform.FromGo(),
			},
			{
				Name:        "missing_context_values",
				Type:        proto.ColumnType_JSON,
				Description: "The missing content values for this policy simulation.",
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "resource_specific_results",
				Type:        proto.ColumnType_JSON,
				Description: "The resource specific results for this policy simulation.",
				Transform:   transform.FromGo(),
			},
			{
				Name:        "organizations_decision_detail",
				Type:        proto.ColumnType_JSON,
				Description: "The organizations decision detail for this policy simulation.",
				Transform:   transform.FromGo(),
			},
			{
				Name:        "permissions_boundary_decision_detail",
				Type:        proto.ColumnType_JSON,
				Description: "The permissions boundary decision detail for this policy simulation.",
				Transform:   transform.FromGo(),
			},
		},
	}
}

type awsIamPolicySimulatorResult struct {
	Action                            string
	Decision                          *string
	DecisionDetails                   map[string]types.PolicyEvaluationDecisionType
	MatchedStatements                 []types.Statement
	MissingContextValues              []string
	OrganizationsDecisionDetail       *types.OrganizationsDecisionDetail
	PermissionsBoundaryDecisionDetail *types.PermissionsBoundaryDecisionDetail
	PrincipalArn                      string
	ResourceArn                       string
	ResourceSpecificResults           []types.ResourceSpecificResult
	Result                            types.EvaluationResult
}

func listIamPolicySimulation(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	principalArn := d.EqualsQuals["principal_arn"].GetStringValue()
	action := d.EqualsQuals["action"].GetStringValue()
	resourceArn := d.EqualsQuals["resource_arn"].GetStringValue()

	// Create Session
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_policy_simulator.listIamPolicySimulation", "client_error", err)
		return nil, err
	}

	var params = &iam.SimulatePrincipalPolicyInput{
		PolicySourceArn: &principalArn,
		ActionNames:     []string{action},
		ResourceArns:    []string{resourceArn},
	}

	op, err := svc.SimulatePrincipalPolicy(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_policy_simulator.listIamPolicySimulation", "api_error", err)
		return nil, err
	}

	evalResults := op.EvaluationResults
	resultForAction := evalResults[0]

	row := awsIamPolicySimulatorResult{
		Action:                            action,
		Decision:                          aws.String(string(resultForAction.EvalDecision)),
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

	if len(resultForAction.MatchedStatements) == 0 {
		row.MatchedStatements = nil
	}
	if len(resultForAction.MissingContextValues) == 0 {
		row.MissingContextValues = nil
	}

	d.StreamListItem(ctx, row)

	return nil, nil
}
