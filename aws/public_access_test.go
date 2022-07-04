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
					"*",
					"arn:aws:iam::111011101110:saml-provider/AWSSSO_DO_NOT_DELETE",
					"arn:aws:iam::111122223333:root",
					"arn:aws:iam::123456789012:user/victor@xyz.com",
					"ecs.amazonaws.com",
					"elasticloadbalancing.amazonaws.com"
				],
				"allowed_principal_account_ids": [
					"*",
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
				"public_access_levels": ["List","Read"],
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
				"access_level": "shared",
				"allowed_organization_ids": [
					"o-123456"
				],
				"allowed_principals": [
					"arn:aws:iam::111011101110:saml-provider/AWSSSO_DO_NOT_DELETE",
					"arn:aws:iam::111122223333:root",
					"arn:aws:iam::123456789012:user/victor@xyz.com",
					"ecs.amazonaws.com",
					"elasticloadbalancing.amazonaws.com"
				],
				"allowed_principal_account_ids": [
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
				"public_access_levels": null,
				"public_statement_ids": null
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
				"access_level": "shared",
				"allowed_organization_ids": [
					"o-123456"
				],
				"allowed_principals": null,
				"allowed_principal_account_ids": null,
				"allowed_principal_federated_identities": null,
				"allowed_principal_services": null,
				"is_public": false,
				"public_access_levels": null,
				"public_statement_ids": null
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
				"access_level": "private",
				"allowed_organization_ids": null,
				"allowed_principals": [
					"arn:aws:cloudwatch:us-east-2:111122223333:alarm:*"
				],
				"allowed_principal_account_ids": null,
				"allowed_principal_federated_identities": null,
				"allowed_principal_services": null,
				"is_public": false,
				"public_access_levels": null,
				"public_statement_ids": null
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
						"Resource": "arn:aws:sns:us-east-2:111122223333:MyTopic",
						"Condition": {
							"StringEquals": {
								"aws:SourceOwner": "444455556666"
							}
						}
					}
				]
			}`,
			`{
				"access_level": "shared",
				"allowed_organization_ids": null,
				"allowed_principals": [
					"ses.amazonaws.com"
				],
				"allowed_principal_account_ids": [
					"444455556666"
				],
				"allowed_principal_federated_identities": null,
				"allowed_principal_services": [
					"ses.amazonaws.com"
				],
				"is_public": false,
				"public_access_levels": null,
				"public_statement_ids": null
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
						"Resource": "arn:aws:sns:us-east-1:111122223333:cloudwatch-alarms",
						"Condition": {
							"StringEquals": {
								"aws:PrincipalAccount": "999988887777"
							}
						}
					}
				]
			}`,
			`{
				"access_level": "shared",
				"allowed_organization_ids": null,
				"allowed_principals": null,
				"allowed_principal_account_ids": [
					"999988887777"
				],
				"allowed_principal_federated_identities": null,
				"allowed_principal_services": null,
				"is_public": false,
				"public_access_levels": null,
				"public_statement_ids": null
			}`,
		},
		{
			`Allow user from same account as the topic account to publish message`,
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
						"Resource": "arn:aws:sns:us-east-1:111122223333:cloudwatch-alarms",
						"Condition": {
							"ArnLike": {
	             	"aws:SourceArn": "arn:aws:iam::111122223333:users/*"
	           	}
						}
					}
				]
			}`,
			`{
				"access_level": "private",
				"allowed_organization_ids": null,
				"allowed_principals": [
					"arn:aws:iam::111122223333:users/*"
				],
				"allowed_principal_account_ids": null,
				"allowed_principal_federated_identities": null,
				"allowed_principal_services": null,
				"is_public": false,
				"public_access_levels": null,
				"public_statement_ids": null
			}`,
		},
		{
			`Allow alarms from specific account to publish message to topic`,
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
						"Resource": "arn:aws:sns:us-east-1:123456789012:cloudwatch-alarms",
						"Condition": {
							"StringEquals": {
								"aws:PrincipalAccount": "999988887777"
							},
							"ArnLike": {
	             	"aws:SourceArn": "arn:aws:cloudwatch:us-east-1:*:alarm:*"
	           	}
						}
					}
				]
			}`,
			`{
        	"access_level": "shared",
        	"allowed_organization_ids": null,
        	"allowed_principals": [
        		"arn:aws:cloudwatch:us-east-1:*:alarm:*"
        	],
        	"allowed_principal_account_ids": [
        		"999988887777"
        	],
        	"allowed_principal_federated_identities": null,
        	"allowed_principal_services": null,
        	"is_public": false,
        	"public_access_levels": null,
        	"public_statement_ids": null
        }`,
		},
		{
			`Doesn't allow user from account 999988887777 to publish to a topic that is owned by another account 111122223333`,
			9,
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
						"Resource": "arn:aws:sns:us-east-1:111122223333:MyTopic",
						"Condition": {
							"ArnLike": {
	             	"aws:SourceArn": "arn:aws:cloudwatch:us-east-1:123456789012:alarm:*"
	           	}
						}
					}
				]
			}`,
			`{
				"access_level": "shared",
				"allowed_organization_ids": null,
				"allowed_principals": [
					"arn:aws:cloudwatch:us-east-1:123456789012:alarm:*"
				],
				"allowed_principal_account_ids": [
					"123456789012"
				],
				"allowed_principal_federated_identities": null,
				"allowed_principal_services": null,
				"is_public": false,
				"public_access_levels": null,
				"public_statement_ids": null
			}`,
		},
		{
			`Private access with the use of aws:SourceIp condition key`,
			10,
			`{
				"Statement": [
					{
						"Action": [
							"s3:GetBucketLocation",
							"s3:ListBucket"
						],
						"Condition": {
							"ForAnyValue:IpAddress": {
								"aws:SourceIp": "122.161.78.130"
							}
						},
						"Effect": "Allow",
						"Principal": {
							"AWS": "*"
						},
						"Resource": "arn:aws:s3:::osborn-shaktiman-bucket-share",
						"Sid": "Example permissions"
					}
				],
				"Version": "2012-10-17"
			}`,
			`{
				"access_level": "shared",
				"allowed_organization_ids": null,
				"allowed_principals": [
					"*"
				],
				"allowed_principal_account_ids": null,
				"allowed_principal_federated_identities": null,
				"allowed_principal_services": null,
				"is_public": false,
				"public_access_levels": null,
				"public_statement_ids": null
			}`,
		},
		{
			`Shared access with the use of aws:PrincipalArn condition key`,
			11,
			`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Sid": "Statement1",
						"Effect": "Allow",
						"Principal": "*",
						"Action": "s3:ListBucket",
						"Resource": "arn:aws:s3:::test-anonymous-access",
						"Condition": {
							"ForAnyValue:ArnLike": {
								"aws:PrincipalArn": [
									"arn:aws:iam::111122223333:root",
									"arn:aws:iam::111122224444:root"
								]
							}
						}
					}
				]
			}`,
			`{
				"access_level": "shared",
				"allowed_organization_ids": null,
				"allowed_principals": [
					"arn:aws:iam::111122223333:root",
					"arn:aws:iam::111122224444:root"
				],
				"allowed_principal_account_ids": [
					"111122224444"
				],
				"allowed_principal_federated_identities": null,
				"allowed_principal_services": null,
				"is_public": false,
				"public_access_levels": null,
				"public_statement_ids": null
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
			evaluatedObj, err := policyObject.EvaluatePolicy("111122223333")
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

			if !reflect.DeepEqual(&expectedObj, evaluatedObj) {
				strdata, _ := json.MarshalIndent(evaluatedObj, "", "\t")
				data := new(bytes.Buffer)
				_ = json.Indent(data, []byte(test.expected), "", "\t")
				t.Errorf("FAILED: \nExpected:\n %v\n\nEvaluated \n%v\n", data.String(), string(strdata))
			}
		})
	}
}

func TestS3ResourcePublicPolicies(t *testing.T) {
	testCases := []publicAccessTest{
		{
			`* principal public access`,
			1,
			`{
				"Statement": [
					{
						"Action": "s3:ListBucket",
						"Effect": "Allow",
						"Principal": "*",
						"Resource": "arn:aws:s3:::test-anonymous-access",
						"Sid": "Statement1"
					}
				],
				"Version": "2012-10-17"
			}`,
			`{
				"access_level": "public",
				"allowed_organization_ids": null,
				"allowed_principals": [
					"*"
				],
				"allowed_principal_account_ids": [
					"*"
				],
				"allowed_principal_federated_identities": null,
				"allowed_principal_services": null,
				"is_public": true,
				"public_access_levels": ["List"],
				"public_statement_ids": [
					"Statement1"
				]
			}`,
		},
		{
			`* principal public access with StringEqualsIfExists on aws:PrincipalAccount`,
			2,
			`{
				"Statement": [
					{
					"Action": "s3:ListBucket",
					"Condition": {
						"StringEqualsIfExists": {
						"aws:PrincipalAccount": "111122223333"
						}
					},
					"Effect": "Allow",
					"Principal": "*",
					"Resource": "arn:aws:s3:::test-anonymous-access",
					"Sid": "Statement1"
					}
				],
				"Version": "2012-10-17"
			}`,
			`{
				"access_level": "public",
				"allowed_organization_ids": null,
				"allowed_principals": [
					"*"
				],
				"allowed_principal_account_ids": [
					"*"
				],
				"allowed_principal_federated_identities": null,
				"allowed_principal_services": null,
				"is_public": true,
				"public_access_levels": ["List"],
				"public_statement_ids": [
					"Statement1"
				]
			}`,
		},
		{
			`* principal public access with StringEqualsIfExists on aws:PrincipalAccount and StringEquals on aws:username`,
			3,
			`{
				"Statement": [
					{
						"Action": "s3:ListBucket",
						"Condition": {
							"StringEquals": {
								"aws:username": "lalit"
							},
							"StringEqualsIfExists": {
								"aws:PrincipalAccount": "111122223333"
							}
						},
						"Effect": "Allow",
						"Principal": "*",
						"Resource": "arn:aws:s3:::test-anonymous-access",
						"Sid": "Statement1"
					}
				],
				"Version": "2012-10-17"
			}`,
			`{
				"access_level": "public",
				"allowed_organization_ids": null,
				"allowed_principals": [
					"*"
				],
				"allowed_principal_account_ids": [
					"*"
				],
				"allowed_principal_federated_identities": null,
				"allowed_principal_services": null,
				"is_public": true,
				"public_access_levels": ["List"],
				"public_statement_ids": [
					"Statement1"
				]
			}`,
		},
		{
			`* principal public access with ArnLike on aws:PrincipalArn with arn for all iam users`,
			4,
			`{
				"Statement": [
					{
						"Action": "s3:ListBucket",
						"Condition": {
							"ArnLike": {
								"aws:PrincipalArn": "arn:aws:iam::*:user/*"
							}
						},
						"Effect": "Allow",
						"Principal": "*",
						"Resource": "arn:aws:s3:::test-anonymous-access",
						"Sid": "Statement1"
					}
				],
				"Version": "2012-10-17"
			}`,
			`{
				"access_level": "public",
				"allowed_organization_ids": null,
				"allowed_principals": [
					"*",
					"arn:aws:iam::*:user/*"
				],
				"allowed_principal_account_ids": [
					"*"
				],
				"allowed_principal_federated_identities": null,
				"allowed_principal_services": null,
				"is_public": true,
				"public_access_levels": ["List"],
				"public_statement_ids": [
					"Statement1"
				]
			}`,
		},
		{
			`* principal public access with ArnLike on aws:SourceArn with arn for all iam services`,
			5,
			`{
				"Statement": [
					{
						"Action": "s3:ListBucket",
						"Condition": {
							"ArnLike": {
								"aws:SourceArn": "arn:aws:iam::*:*/*"
							}
						},
						"Effect": "Allow",
						"Principal": "*",
						"Resource": "arn:aws:s3:::test-anonymous-access",
						"Sid": "Statement1"
					}
				],
				"Version": "2012-10-17"
			}`,
			`{
        	"access_level": "public",
        	"allowed_organization_ids": null,
        	"allowed_principals": [
						"*",
        		"arn:aws:iam::*:*/*"
        	],
        	"allowed_principal_account_ids": [
        		"*"
        	],
        	"allowed_principal_federated_identities": null,
        	"allowed_principal_services": null,
        	"is_public": true,
        	"public_access_levels": ["List"],
        	"public_statement_ids": [
        		"Statement1"
        	]
        }`,
		},
		{
			`* principal public access with ArnLike on aws:SourceArn with arn for all cloudwatch alarms`,
			6,
			`{
				"Statement": [
					{
						"Action": "s3:ListBucket",
						"Condition": {
							"ArnLike": {
								"aws:SourceArn": "arn:aws:cloudwatch:us-east-1:*:alarm:*"
							}
						},
						"Effect": "Allow",
						"Principal": "*",
						"Resource": "arn:aws:s3:::test-anonymous-access",
						"Sid": "Statement1"
					}
				],
				"Version": "2012-10-17"
			}`,
			`{
					"access_level": "public",
					"allowed_organization_ids": null,
					"allowed_principals": [
						"*",
						"arn:aws:cloudwatch:us-east-1:*:alarm:*"
					],
					"allowed_principal_account_ids": [
						"*"
					],
					"allowed_principal_federated_identities": null,
					"allowed_principal_services": null,
					"is_public": true,
					"public_access_levels": ["List"],
					"public_statement_ids": [
						"Statement1"
					]
				}`,
		},
		{
			`* principal public access with ArnLike on aws:PrincipalArn with arn`,
			7,
			`{
				"Statement": [
					{
						"Action": "s3:ListBucket",
						"Condition": {
							"ForAnyValue:ArnLike": {
								"aws:PrincipalArn": [
									"arn:aws:iam::*:root",
									"arn:aws:iam::444422223333:root"
								]
							}
						},
						"Effect": "Allow",
						"Principal": "*",
						"Resource": "arn:aws:s3:::test-anonymous-access",
						"Sid": "Statement1"
					}
				],
				"Version": "2012-10-17"
			}`,
			`{
				"access_level": "public",
				"allowed_organization_ids": null,
				"allowed_principals": [
					"*",
					"arn:aws:iam::*:root",
					"arn:aws:iam::444422223333:root"
				],
				"allowed_principal_account_ids": [
					"*",
					"444422223333"
				],
				"allowed_principal_federated_identities": null,
				"allowed_principal_services": null,
				"is_public": true,
				"public_access_levels": ["List"],
				"public_statement_ids": [
					"Statement1"
				]
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
			evaluatedObj, err := policyObject.EvaluatePolicy("111122223333")
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

			if !reflect.DeepEqual(&expectedObj, evaluatedObj) {
				strdata, _ := json.MarshalIndent(evaluatedObj, "", "\t")
				data := new(bytes.Buffer)
				_ = json.Indent(data, []byte(test.expected), "", "\t")
				t.Errorf("FAILED: \nExpected:\n %v\n\nEvaluated \n%v\n", data.String(), string(strdata))
			}
		})
	}
}

// https://docs.aws.amazon.com/AmazonS3/latest/userguide/example-bucket-policies.html
func TestS3ExampleResourcePolicies(t *testing.T) {
	testCases := []publicAccessTest{
		{
			`Granting permissions to multiple accounts with added conditions`,
			1,
			`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Sid": "AddCannedAcl",
						"Effect": "Allow",
						"Principal": {
							"AWS": [
								"arn:aws:iam::111122223333:root",
								"arn:aws:iam::444455556666:root"
							]
						},
						"Action": [
							"s3:PutObject",
							"s3:PutObjectAcl"
						],
						"Resource": "arn:aws:s3:::DOC-EXAMPLE-BUCKET/*",
						"Condition": {
							"StringEquals": {
								"s3:x-amz-acl": [
									"public-read"
								]
							}
						}
					}
				]
			}`,
			`{
				"access_level": "shared",
				"allowed_organization_ids": [],
				"allowed_principals": [
					"arn:aws:iam::111122223333:root",
					"arn:aws:iam::444455556666:root"
				],
				"allowed_principal_account_ids": [
					"111122223333",
					"444455556666"
				],
				"allowed_principal_federated_identities": [],
				"allowed_principal_services": [],
				"is_public": false,
				"public_access_levels": ["Permissions management","Write"],
				"public_statement_ids": []
			}`,
		},
		{
			`Granting read-only permission to an anonymous user`,
			2,
			`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Sid": "PublicRead",
						"Effect": "Allow",
						"Principal": "*",
						"Action": [
							"s3:GetObject",
							"s3:GetObjectVersion"
						],
						"Resource": [
							"arn:aws:s3:::DOC-EXAMPLE-BUCKET/*"
						]
					}
				]
			}`,
			`{
				"access_level": "public",
				"allowed_organization_ids": [],
				"allowed_principals": [
					"*"
				],
				"allowed_principal_account_ids": [
					"*"
				],
				"allowed_principal_federated_identities": [],
				"allowed_principal_services": [],
				"is_public": true,
				"public_access_levels": ["Read"],
				"public_statement_ids": [
					"PublicRead"
				]
			}`,
		},
		{
			`Limiting access to specific IP addresses`,
			3,
			`{
				"Version": "2012-10-17",
				"Id": "S3PolicyId1",
				"Statement": [
					{
						"Sid": "IPAllow",
						"Effect": "Deny",
						"Principal": "*",
						"Action": "s3:*",
						"Resource": [
							"arn:aws:s3:::DOC-EXAMPLE-BUCKET",
							"arn:aws:s3:::DOC-EXAMPLE-BUCKET/*"
						],
						"Condition": {
							"NotIpAddress": {
								"aws:SourceIp": "54.240.143.0/24"
							}
						}
					}
				]
			}`,
			`{
        	"access_level": "private",
        	"allowed_organization_ids": [],
        	"allowed_principals": [],
        	"allowed_principal_account_ids": [],
        	"allowed_principal_federated_identities": [],
        	"allowed_principal_services": [],
        	"is_public": false,
        	"public_access_levels": [],
        	"public_statement_ids": []
        }`,
		},
		{
			`Policy allows request with mfa only`,
			4,
			`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Sid": "Statement1",
						"Effect": "Allow",
						"Principal": "*",
						"Action": [
							"s3:ListBucket",
							"s3:GetObject"
						],
						"Resource": [
							"arn:aws:s3:::test-anonymous-access",
							"arn:aws:s3:::test-anonymous-access/*"
						],
						"Condition": { "Null": { "aws:MultiFactorAuthAge": true }}
					}
				]
			}`,
			`{
        	"access_level": "public",
        	"allowed_organization_ids": [],
        	"allowed_principals": [
        		"*"
        	],
        	"allowed_principal_account_ids": [
        		"*"
        	],
        	"allowed_principal_federated_identities": [],
        	"allowed_principal_services": [],
        	"is_public": true,
        	"public_access_levels": ["List","Read"],
        	"public_statement_ids": [
        		"Statement1"
        	]
        }`,
		},
		{
			`Policy allows request with mfa and principal account limitation`,
			5,
			`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Sid": "Statement1",
						"Effect": "Allow",
						"Principal": "*",
						"Action": [
							"s3:ListBucket",
							"s3:GetObject"
						],
						"Resource": [
							"arn:aws:s3:::test-anonymous-access",
							"arn:aws:s3:::test-anonymous-access/*"
						],
						"Condition": {
							"Null": {
								"aws:MultiFactorAuthAge": true
							},
							"StringEquals": {
								"aws:PrincipalAccount": "999988887777"
							}
						}
					}
				]
			}`,
			`{
        	"access_level": "shared",
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
        	"public_access_levels": ["List", "Read"],
        	"public_statement_ids": []
        }`,
		},
		{
			`Granting cross-account permissions to upload objects while ensuring the bucket owner has full control`,
			6,
			`{
				"Version":"2012-10-17",
				"Statement":[
					{
						"Sid":"PolicyForAllowUploadWithACL",
						"Effect":"Allow",
						"Principal":{"AWS":"111122223333"},
						"Action":"s3:PutObject",
						"Resource":"arn:aws:s3:::DOC-EXAMPLE-BUCKET/*",
						"Condition": {
							"StringEquals": {"s3:x-amz-acl":"bucket-owner-full-control"}
						}
					}
				]
			}`,
			`{
				"access_level": "shared",
				"allowed_organization_ids": [],
				"allowed_principals": [
					"111122223333"
				],
				"allowed_principal_account_ids": [
					"111122223333"
				],
				"allowed_principal_federated_identities": [],
				"allowed_principal_services": [],
				"is_public": false,
				"public_access_levels": ["Write"],
				"public_statement_ids": []
			}`,
		},
		{
			`Granting permissions for Amazon S3 Inventory and Amazon S3 analytics`,
			7,
			`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Sid": "InventoryAndAnalyticsExamplePolicy",
						"Effect": "Allow",
						"Principal": {
							"Service": "s3.amazonaws.com"
						},
						"Action": "s3:PutObject",
						"Resource": [
							"arn:aws:s3:::destinationbucket/*"
						],
						"Condition": {
							"ArnLike": {
								"aws:SourceArn": "arn:aws:s3:::sourcebucket"
							},
							"StringEquals": {
								"aws:SourceAccount": "111122223333",
								"s3:x-amz-acl": "bucket-owner-full-control"
							}
						}
					}
				]
			}`,
			`{
				"access_level": "shared",
				"allowed_organization_ids": [],
				"allowed_principal_account_ids": [
					"111122223333"
				],
				"allowed_principal_federated_identities": [],
				"allowed_principal_services": [
					"s3.amazonaws.com"
				],
				"allowed_principals": [
					"111122223333",
					"arn:aws:s3:::sourcebucket",
					"s3.amazonaws.com"
				],
				"is_public": false,
				"public_access_levels": [
					"Write"
				],
				"public_statement_ids": []
			}`,
		},
		{
			`Granting permissions for Amazon S3 Inventory and Amazon S3 analytics with only source arn condition`,
			8,
			`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Sid": "InventoryAndAnalyticsExamplePolicy",
						"Effect": "Allow",
						"Principal": {
							"Service": "s3.amazonaws.com"
						},
						"Action": "s3:PutObject",
						"Resource": [
							"arn:aws:s3:::destinationbucket/*"
						],
						"Condition": {
							"ArnLike": {
								"aws:SourceArn": "arn:aws:s3:::sourcebucket"
							}
						}
					}
				]
			}`,
			`{
				"access_level": "private",
				"allowed_organization_ids": [],
				"allowed_principals": [
					"arn:aws:s3:::sourcebucket",
					"s3.amazonaws.com"
				],
				"allowed_principal_account_ids": [],
				"allowed_principal_federated_identities": [],
				"allowed_principal_services": [
					"s3.amazonaws.com"
				],
				"is_public": false,
				"public_access_levels": [
					"Write"
				],
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
			evaluatedObj, err := policyObject.EvaluatePolicy("111122223333")
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

			if !reflect.DeepEqual(&expectedObj, evaluatedObj) {
				strdata, _ := json.MarshalIndent(evaluatedObj, "", "\t")
				data := new(bytes.Buffer)
				_ = json.Indent(data, []byte(test.expected), "", "\t")
				t.Errorf("FAILED: \nExpected:\n %v\n\nEvaluated \n%v\n", data.String(), string(strdata))
			}
		})
	}
}

/*
func TestAccessPoliciesActions(t *testing.T) {
	permissionsData := getParliamentIamPermissions()
	testCases := []struct {
		Actions      []string
		AccessLevels []string
	}{
		{
			[]string{"s3:PutObject", "s3:PutObjectAcl", "s3:list*"},
			[]string{"List", "Permissions management", "Write"},
		},
		{
			[]string{"*"},
			[]string{"List", "Permissions management", "Read", "Tagging", "Write"},
		},
		{
			[]string{"s3:*"},
			[]string{"List", "Permissions management", "Read", "Tagging", "Write"},
		},
		{
			[]string{"iam:*"},
			[]string{"List", "Permissions management", "Read", "Tagging", "Write"},
		},
		{
			[]string{"iam:create*"},
			[]string{"Permissions management", "Write"},
		},
		{
			[]string{"iam:create*"},
			[]string{"Permissions management", "Write"},
		},
		{
			[]string{"iam:create*"},
			[]string{"Permissions management", "Write"},
		},
		{
			[]string{
				"SNS:GetTopicAttributes",
				"SNS:SetTopicAttributes",
				"SNS:AddPermission",
				"SNS:RemovePermission",
				"SNS:DeleteTopic",
				"SNS:Subscribe",
				"SNS:ListSubscriptionsByTopic",
				"SNS:Publish",
				"SNS:Receive",
			},
			[]string{"List", "Permissions management", "Read", "Write"},
		},
	}

	for i, test := range testCases {
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			fmt.Println("start", i+1, time.Now())
			accessLevels := GetAccessLevelsFromActions(permissionsData, test.Actions)
			fmt.Println("end", i+1, time.Now())

			// Sort []string attributes to compare
			sort.Strings(accessLevels)
			sort.Strings(test.AccessLevels)
			// fmt.Println(accessLevels)

			if !reflect.DeepEqual(accessLevels, test.AccessLevels) {
				t.Errorf("FAILED: \nExpected:\n %v\n\nGot:\n %v\n", test.AccessLevels, accessLevels)
			}
		})
	}
}
*/
