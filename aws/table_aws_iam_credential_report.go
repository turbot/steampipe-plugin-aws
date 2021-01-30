package aws

import (
	"context"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//type awsIamPolicySimulatorResult struct {
//	PrincipalArn string
//	Action       string
//	ResourceArn  string
//	Decision     *string
//	Result       *iam.EvaluationResult
//}

func tableAwsIamCredentialReport(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_credential_report",
		Description: "AWS IAM Credential Report",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("generated_time"),
			Hydrate:    getIamCredentialReport,
		},
		Columns: []*plugin.Column{
			// "Key" Columns
			{
				Name:        "generated_time",
				Description: "The date and time when the credential report was created, in ISO 8601 date-time format (http://www.iso.org/iso/iso8601)",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "report_format",
				Description: "The format (MIME type) of the credential report.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
		},
	}
}

func getIamCredentialReport(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getIamCredentialReport")
	generatedTime := d.KeyColumnQuals["generated_time"].GetStringValue()
	action := d.KeyColumnQuals["action"].GetStringValue()
	resourceArn := d.KeyColumnQuals["resource_arn"].GetStringValue()

	// Create Session
	svc, err := IAMService(ctx, d.ConnectionManager)
	if err != nil {
		return nil, err
	}

	var params = &iam.SimulatePrincipalPolicyInput{
		PolicySourceArn: &generatedTime,
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
		PrincipalArn: generatedTime,
		Action:       action,
		ResourceArn:  resourceArn,
		Decision:     resultForAction.EvalDecision,
		Result:       resultForAction,
	}

	return row, nil
}
