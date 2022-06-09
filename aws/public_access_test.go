package aws

import (
	"encoding/json"
	"fmt"
	"testing"
)

type numberConversionTest struct {
	name     string
	input    string
	expected string
}

// go test -v -run ^TestGetConditionalKeymapping$ github.com/turbot/steampipe-plugin-aws/aws

func TestGetConditionalKeymapping(t *testing.T) {
	testCases := []numberConversionTest{
		{
			`any_value`,
			`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":"*"},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":"arn:aws:s3:::test"},{"Sid":"OrganizationAccess","Effect":"Allow","Principal":{"AWS":"*"},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":"arn:aws:s3:::test","Condition":{"StringEquals":{"aws:PrincipalOrgID":["o-123456"]}}},{"Sid":"AccountPrincipals","Effect":"Allow","Principal":{"AWS":["arn:aws:iam::123456789012:user/victor@xyz.com","arn:aws:iam::111122223333:root"]},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":"arn:aws:s3:::test"},{"Sid":"FederatedPrincipals","Effect":"Allow","Principal":{"Federated":"arn:aws:iam::111011101110:saml-provider/AWSSSO_DO_NOT_DELETE"},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":"arn:aws:s3:::test"},{"Sid":"ServicePrincipals","Effect":"Allow","Principal":{"Service":["ecs.amazonaws.com","elasticloadbalancing.amazonaws.com"]},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":"arn:aws:s3:::test"},{"Sid":"PublicAccess","Effect":"Allow","Principal":{"AWS":"*"},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":"arn:aws:s3:::test"}]}`,
			``,
		},
		// {
		// 	`single_value`,
		// 	`{"Version":"2012-10-17","Statement":[{"Sid":"","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::111122223333:root"},"Action":"sts:AssumeRole","Condition":{"StringEquals":{"aws:PrincipalAccount":"123456789012","aws:PrincipalOrgID":"o-123456789"}}}]}`,
		// 	`{"aws:principalaccount":{"StringEquals":{"Key":"aws:principalaccount","Value":["123456789012"]}},"aws:principalorgid":{"StringEquals":{"Key":"aws:principalorgid","Value":["o-123456789"]}}}`,
		// },
	}

	for _, test := range testCases {

		policy, err := canonicalPolicy(test.input)
		if err != nil {
			t.Errorf("Convert failed for case '%s': %v", test.input, err)
		}

		policyObject, ok := policy.(Policy)
		if !ok {
			t.Errorf("Unable to parse input as policy")
		}

		evaluation, err := policyObject.EvaluatePolicy()
		if err != nil {
			t.Errorf("Unable to parse input as policy")
		}

		// var input Policy
		// _ = json.Unmarshal([]byte(test.expected), &input)
		strdata, _ := json.MarshalIndent(evaluation, "", "\t")
		// output, _ := json.MarshalIndent(newCondition, "", "\t")
		fmt.Printf("%s\n", string(strdata))

		// if !reflect.DeepEqual(input, *newCondition) {
		// 	t.Errorf("\nTest: '%s.%s' FAILED\nexpected:\n %v \ngot:\n %v \n", "TestConvertStatementCondition", test.name, input, *newCondition)
		// }
		// fmt.Printf("\nTest: '%s.%s' PASSED\noutput:\n %v\n", "TestConvertStatementCondition", test.name, input)
	}
}
