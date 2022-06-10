package aws

import (
	"encoding/json"
	"fmt"
	"testing"
)

type publicAccessTest struct {
	name     string
	policy   string
	expected string
}

// go test -v -run ^TestGetConditionalKeymapping$ github.com/turbot/steampipe-plugin-aws/aws

func TestGetConditionalKeymapping(t *testing.T) {
	testCases := []publicAccessTest{
		{
			`1. AWS S3 Multiple statements with public access`,
			`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Effect": "Allow",
						"Principal": {
							"AWS": "*"
						},
						"Action": [
							"s3:GetBucketLocation",
							"s3:ListBucket"
						],
						"Resource": "arn:aws:s3:::test"
					},
					{
						"Sid": "OrganizationAccess",
						"Effect": "Allow",
						"Principal": {
							"AWS": "*"
						},
						"Action": [
							"s3:GetBucketLocation",
							"s3:ListBucket"
						],
						"Resource": "arn:aws:s3:::test",
						"Condition": {
							"StringEquals": {
								"aws:PrincipalOrgID": [
									"o-123456"
								]
							}
						}
					},
					{
						"Sid": "AccountPrincipals",
						"Effect": "Allow",
						"Principal": {
							"AWS": [
								"arn:aws:iam::123456789012:user/victor@xyz.com",
								"arn:aws:iam::111122223333:root"
							]
						},
						"Action": [
							"s3:GetBucketLocation",
							"s3:ListBucket"
						],
						"Resource": "arn:aws:s3:::test"
					},
					{
						"Sid": "FederatedPrincipals",
						"Effect": "Allow",
						"Principal": {
							"Federated": "arn:aws:iam::111011101110:saml-provider/AWSSSO_DO_NOT_DELETE"
						},
						"Action": [
							"s3:GetBucketLocation",
							"s3:ListBucket"
						],
						"Resource": "arn:aws:s3:::test"
					},
					{
						"Sid": "ServicePrincipals",
						"Effect": "Allow",
						"Principal": {
							"Service": [
								"ecs.amazonaws.com",
								"elasticloadbalancing.amazonaws.com"
							]
						},
						"Action": [
							"s3:GetBucketLocation",
							"s3:ListBucket"
						],
						"Resource": "arn:aws:s3:::test"
					},
					{
						"Sid": "PublicAccess",
						"Effect": "Allow",
						"Principal": {
							"AWS": "*"
						},
						"Action": [
							"s3:GetBucketLocation",
							"s3:ListBucket"
						],
						"Resource": "arn:aws:s3:::test"
					}
				]
			}`,
			``,
		},
		{
			`2. AWS S3 Multiple statements without public access`,
			`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Sid": "OrganizationAccess",
						"Effect": "Allow",
						"Principal": {
							"AWS": "*"
						},
						"Action": [
							"s3:GetBucketLocation",
							"s3:ListBucket"
						],
						"Resource": "arn:aws:s3:::test",
						"Condition": {
							"StringEquals": {
								"aws:PrincipalOrgID": [
									"o-123456"
								]
							}
						}
					},
					{
						"Sid": "AccountPrincipals",
						"Effect": "Allow",
						"Principal": {
							"AWS": [
								"arn:aws:iam::123456789012:user/victor@xyz.com",
								"arn:aws:iam::111122223333:root"
							]
						},
						"Action": [
							"s3:GetBucketLocation",
							"s3:ListBucket"
						],
						"Resource": "arn:aws:s3:::test"
					},
					{
						"Sid": "FederatedPrincipals",
						"Effect": "Allow",
						"Principal": {
							"Federated": "arn:aws:iam::111011101110:saml-provider/AWSSSO_DO_NOT_DELETE"
						},
						"Action": [
							"s3:GetBucketLocation",
							"s3:ListBucket"
						],
						"Resource": "arn:aws:s3:::test"
					},
					{
						"Sid": "ServicePrincipals",
						"Effect": "Allow",
						"Principal": {
							"Service": [
								"ecs.amazonaws.com",
								"elasticloadbalancing.amazonaws.com"
							]
						},
						"Action": [
							"s3:GetBucketLocation",
							"s3:ListBucket"
						],
						"Resource": "arn:aws:s3:::test"
					}
				]
			}`,
			``,
		},
		{
			`3. single_statement`,
			`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Sid": "OrganizationAccess",
						"Effect": "Allow",
						"Principal": {
							"AWS": "*"
						},
						"Action": [
							"s3:GetBucketLocation",
							"s3:ListBucket"
						],
						"Resource": "arn:aws:s3:::test",
						"Condition": {
							"StringEquals": {
								"aws:PrincipalOrgID": [
									"o-123456"
								]
							}
						}
					}
				]
			}`,
			``,
		},
		{
			`4. single_statement_with_source_arn`,
			`{
				"Statement": [
					{
						"Effect": "Allow",
						"Principal": {
							"AWS": "*"
						},
						"Action": "SNS:Publish",
						"Resource": "arn:aws:sns:us-east-2:444455556666:MyTopic",
						"Condition": {
							"ArnLike": {
								"aws:SourceArn": "arn:aws:cloudwatch:us-east-2:111122223333:alarm:*"
							}
						}
					}
				]
			}`,
			``,
		},
		{
			`5. Allow Amazon SES to publish to a topic that is owned by another account`,
			`{
				"Version": "2008-10-17",
				"Id": "__default_policy_ID",
				"Statement": [
					{
						"Sid": "__default_statement_ID",
						"Effect": "Allow",
						"Principal": {
							"Service": "ses.amazonaws.com"
						},
						"Action": "SNS:Publish",
						"Resource": "arn:aws:sns:us-east-2:444455556666:MyTopic",
						"Condition": {
							"StringEquals": {
								"aws:SourceOwner": "111122223333"
							}
						}
					}
				]
			}`,
			``,
		},
		{
			`6. Allow user from account 999988887777 to publish to a topic that is owned by another account 123456789012`,
			`{
				"Version": "2008-10-17",
				"Id": "__default_policy_ID",
				"Statement": [
					{
						"Sid": "__default_statement_ID",
						"Effect": "Allow",
						"Principal": {
							"AWS": "*"
						},
						"Action": [
							"SNS:GetTopicAttributes",
							"SNS:SetTopicAttributes",
							"SNS:AddPermission",
							"SNS:RemovePermission",
							"SNS:DeleteTopic",
							"SNS:Subscribe",
							"SNS:ListSubscriptionsByTopic",
							"SNS:Publish",
							"SNS:Receive"
						],
						"Resource": "arn:aws:sns:us-east-1:123456789012:cloudwatch-alarms",
						"Condition": {
							"StringEquals": {
								"aws:PrincipalAccount": "999988887777"
							}
						}
					}
				]
			}`,
			``,
		},
		{
			`7. Allow user from same account as the topic accout to publish message`,
			`{
				"Version": "2008-10-17",
				"Id": "__default_policy_ID",
				"Statement": [
					{
						"Sid": "__default_statement_ID",
						"Effect": "Allow",
						"Principal": {
							"AWS": "*"
						},
						"Action": [
							"SNS:GetTopicAttributes",
							"SNS:SetTopicAttributes",
							"SNS:AddPermission",
							"SNS:RemovePermission",
							"SNS:DeleteTopic",
							"SNS:Subscribe",
							"SNS:ListSubscriptionsByTopic",
							"SNS:Publish",
							"SNS:Receive"
						],
						"Resource": "arn:aws:sns:us-east-1:123456789012:cloudwatch-alarms",
						"Condition": {
							"StringEquals": {
								"aws:PrincipalAccount": "999988887777"
							},
							"ArnLike": {
	             	"aws:SourceArn": "arn:aws:cloudwatch:us-east-1:111122223333:alarm:*"
	           	}
						}
					}
				]
			}`,
			``,
		},
		{
			`8. Doesn't allow user from account 999988887777 to publish to a topic that is owned by another account 123456789012`,
			`{
				"Version": "2008-10-17",
				"Id": "__default_policy_ID",
				"Statement": [
					{
						"Sid": "__default_statement_ID",
						"Effect": "Allow",
						"Principal": {
							"AWS": "*"
						},
						"Action": [
							"SNS:GetTopicAttributes",
							"SNS:SetTopicAttributes",
							"SNS:AddPermission",
							"SNS:RemovePermission",
							"SNS:DeleteTopic",
							"SNS:Subscribe",
							"SNS:ListSubscriptionsByTopic",
							"SNS:Publish",
							"SNS:Receive"
						],
						"Resource": "arn:aws:sns:us-east-1:123456789012:cloudwatch-alarms",
						"Condition": {
							"StringEquals": {
								"aws:PrincipalAccount": "999988887777"
							},
							"ArnLike": {
	             	"aws:SourceArn": "arn:aws:cloudwatch:us-east-1:111122223333:alarm:*"
	           	}
						}
					}
				]
			}`,
			``,
		},
	}

	for _, test := range testCases {

		policy, err := canonicalPolicy(test.policy)
		if err != nil {
			t.Errorf("Convert failed for case '%s': %v", test.policy, err)
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
		fmt.Printf("\n%s:\n%s\n", test.name, string(strdata))

		// if !reflect.DeepEqual(input, *newCondition) {
		// 	t.Errorf("\nTest: '%s.%s' FAILED\nexpected:\n %v \ngot:\n %v \n", "TestConvertStatementCondition", test.name, input, *newCondition)
		// }
		// fmt.Printf("\nTest: '%s.%s' PASSED\noutput:\n %v\n", "TestConvertStatementCondition", test.name, input)
	}
}
