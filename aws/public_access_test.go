package aws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"testing"
)

type publicAccessTest struct {
	name     string
	testNo   int
	policy   string
	expected string
}

// go test -v -run ^TestResourcePolicyPublicAccess$ github.com/turbot/steampipe-plugin-aws/aws
// Run all tests
// go test -v -run TestResourcePolicyPublicAccess/Test

// Run individual test case - I will run test with testNo 1
// go test -v -run TestResourcePolicyPublicAccess/Test=1

func TestResourcePolicyPublicAccess(t *testing.T) {
	testCases := []publicAccessTest{
		{
			`AWS S3 Multiple statements with public access`,
			1,
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
			`{
				"access_level": "public",
				"allowed_organization_ids": [
					"o-123456"
				],
				"allowed_principals": [
					"ecs.amazonaws.com",
					"elasticloadbalancing.amazonaws.com",
					"arn:aws:iam::111011101110:saml-provider/AWSSSO_DO_NOT_DELETE",
					"*",
					"o-123456",
					"arn:aws:iam::111122223333:root",
					"arn:aws:iam::123456789012:user/victor@xyz.com"
				],
				"allowed_principal_account_ids": [
					"*",
					"o-123456",
					"111122223333",
					"123456789012"
				],
				"allowed_principal_federated_identities": [
					"arn:aws:iam::111011101110:saml-provider/AWSSSO_DO_NOT_DELETE"
				],
				"allowed_principal_services": [
					"ecs.amazonaws.com",
					"elasticloadbalancing.amazonaws.com"
				],
				"is_public": true,
				"public_access_levels": [],
				"public_statement_ids": [
					"PublicAccess",
					"Statement[1]"
				]
			}`,
		},
		{
			`AWS S3 Multiple statements without public access`,
			2,
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
			`{
				"access_level": "",
				"allowed_organization_ids": [
					"o-123456"
				],
				"allowed_principals": [
					"arn:aws:iam::123456789012:user/victor@xyz.com",
					"ecs.amazonaws.com",
					"elasticloadbalancing.amazonaws.com",
					"arn:aws:iam::111011101110:saml-provider/AWSSSO_DO_NOT_DELETE",
					"o-123456",
					"arn:aws:iam::111122223333:root"
				],
				"allowed_principal_account_ids": [
					"o-123456",
					"111122223333",
					"123456789012"
				],
				"allowed_principal_federated_identities": [
					"arn:aws:iam::111011101110:saml-provider/AWSSSO_DO_NOT_DELETE"
				],
				"allowed_principal_services": [
					"ecs.amazonaws.com",
					"elasticloadbalancing.amazonaws.com"
				],
				"is_public": false,
				"public_access_levels": [],
				"public_statement_ids": []
			}`,
		},
		{
			`single_statement`,
			3,
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
			`{
				"access_level": "",
				"allowed_organization_ids": [
					"o-123456"
				],
				"allowed_principals": [
					"o-123456"
				],
				"allowed_principal_account_ids": [
					"o-123456"
				],
				"allowed_principal_federated_identities": [],
				"allowed_principal_services": [],
				"is_public": false,
				"public_access_levels": [],
				"public_statement_ids": []
			}`,
		},
		{
			`single_statement_with_source_arn`,
			4,
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
			`{
				"access_level": "",
				"allowed_organization_ids": [],
				"allowed_principals": [
					"arn:aws:cloudwatch:us-east-2:111122223333:alarm:*"
				],
				"allowed_principal_account_ids": [
					"111122223333"
				],
				"allowed_principal_federated_identities": [],
				"allowed_principal_services": [],
				"is_public": false,
				"public_access_levels": [],
				"public_statement_ids": []
			}`,
		},
		{
			`Allow Amazon SES to publish to a topic that is owned by another account`,
			5,
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
			`{
				"access_level": "",
				"allowed_organization_ids": [],
				"allowed_principals": [
					"ses.amazonaws.com"
				],
				"allowed_principal_account_ids": [],
				"allowed_principal_federated_identities": [],
				"allowed_principal_services": [
					"ses.amazonaws.com"
				],
				"is_public": false,
				"public_access_levels": [],
				"public_statement_ids": []
			}`,
		},
		{
			`Allow user from account 999988887777 to publish to a topic that is owned by another account 123456789012`,
			6,
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
			`{
				"access_level": "",
				"allowed_organization_ids": [],
				"allowed_principals": [
					"999988887777"
				],
				"allowed_principal_account_ids": [
					"999988887777"
				],
				"allowed_principal_federated_identities": [],
				"allowed_principal_services": [],
				"is_public": false,
				"public_access_levels": [],
				"public_statement_ids": []
			}`,
		},
		{
			`Allow user from same account as the topic accout to publish message`,
			7,
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
			`{
				"access_level": "",
				"allowed_organization_ids": [],
				"allowed_principals": [
					"arn:aws:cloudwatch:us-east-1:111122223333:alarm:*",
					"cloudwatch.amazonaws.com",
					"999988887777"
				],
				"allowed_principal_account_ids": [
					"111122223333",
					"999988887777"
				],
				"allowed_principal_federated_identities": [],
				"allowed_principal_services": ["cloudwatch.amazonaws.com"],
				"is_public": false,
				"public_access_levels": [],
				"public_statement_ids": []
			}`,
		},
		{
			`Doesn't allow user from account 999988887777 to publish to a topic that is owned by another account 123456789012`,
			8,
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
						"Resource": "arn:aws:sns:us-east-1:123456789012:MyTopic",
						"Condition": {
							"ArnLike": {
	             	"aws:SourceArn": "arn:aws:cloudwatch:us-east-1:111122223333:alarm:*"
	           	}
						}
					}
				]
			}`,
			`{
				"access_level": "",
				"allowed_organization_ids": [],
				"allowed_principals": [
					"arn:aws:cloudwatch:us-east-1:111122223333:alarm:*",
					"cloudwatch.amazonaws.com"
				],
				"allowed_principal_account_ids": [
					"111122223333"
				],
				"allowed_principal_federated_identities": [],
				"allowed_principal_services": ["cloudwatch.amazonaws.com"],
				"is_public": false,
				"public_access_levels": [],
				"public_statement_ids": []
			}`,
		},
	}

	for _, test := range testCases {
		t.Run(fmt.Sprintf("Test=%s %s", strconv.Itoa(test.testNo), test.name), func(t *testing.T) {
			policy, err := canonicalPolicy(test.policy)
			if err != nil {
				t.Errorf("Test: %s Policy canonicalization failed with error: %#v\n", test.name, err)
			}

			policyObject, ok := policy.(Policy)
			if !ok {
				t.Errorf("Test: %s Policy coercion failed with error: %#v\n", test.name, err)
			}
			evaluatedObj, err := policyObject.EvaluatePolicy()
			if err != nil {
				t.Errorf("Test: %s\nPolicy evaluation failed with error: %#v\n", test.name, err)
			}

			var expectedObj PolicyEvaluation
			_ = json.Unmarshal([]byte(test.expected), &expectedObj)

			// Sort []string attributes to compare
			sort.Strings(expectedObj.AllowedOrganizationIds)
			sort.Strings(expectedObj.AllowedPrincipalAccountIds)
			sort.Strings(expectedObj.AllowedPrincipalFederatedIdentities)
			sort.Strings(expectedObj.AllowedPrincipalServices)
			sort.Strings(expectedObj.AllowedPrincipals)
			sort.Strings(expectedObj.PublicAccessLevels)
			sort.Strings(expectedObj.PublicStatementIds)
			sort.Strings(evaluatedObj.AllowedOrganizationIds)
			sort.Strings(evaluatedObj.AllowedPrincipalAccountIds)
			sort.Strings(evaluatedObj.AllowedPrincipalFederatedIdentities)
			sort.Strings(evaluatedObj.AllowedPrincipalServices)
			sort.Strings(evaluatedObj.AllowedPrincipals)
			sort.Strings(evaluatedObj.PublicAccessLevels)
			sort.Strings(evaluatedObj.PublicStatementIds)
			if !reflect.DeepEqual(&expectedObj, evaluatedObj) {
				strdata, _ := json.MarshalIndent(evaluatedObj, "", "\t")
				data := new(bytes.Buffer)
				_ = json.Indent(data, []byte(test.expected), "", "\t")
				t.Errorf("FAILED: \nExpected:\n %v\n\nEvaluated \n%v\n", data.String(), string(strdata))
			}
		})
	}
}
