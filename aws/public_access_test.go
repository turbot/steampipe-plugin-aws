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
					"ServicePrincipals",
					"Statement[1]"
				]
			}`,
		},
		{
			`AWS S3 Multiple statements without public access`,
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
        		"arn:aws:iam::123456789012:user/victor@xyz.com"
        	],
        	"allowed_principal_account_ids": [
        		"123456789012"
        	],
        	"allowed_principal_federated_identities": [
        		"arn:aws:iam::111011101110:saml-provider/AWSSSO_DO_NOT_DELETE"
        	],
        	"allowed_principal_services": null,
        	"is_public": false,
        	"public_access_levels": null,
        	"public_statement_ids": null
        }`,
		},
		{
			`single_statement`,
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
			`Allow Amazon SES to publish to a topic that is owned by another account`,
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
			`Allow alarms from specific account to publish message to topic`,
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

	for index, test := range testCases {
		t.Run(strconv.Itoa(index+1), func(t *testing.T) {
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
			`* principal public access with ArnLike on aws:PrincipalArn with arn to allow root user from any account`,
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
		{
			`Bucket logging policies`,
			`{
				"Statement": [
					{
					"Action": "s3:GetBucketAcl",
					"Effect": "Allow",
					"Principal": {
						"Service": "cloudtrail.amazonaws.com"
					},
					"Resource": "arn:aws:s3:::turbot-111122224444-ap-northeast-1",
					"Sid": "AWSCloudTrailAclCheck"
					},
					{
					"Action": "s3:PutObject",
					"Condition": {
						"StringEquals": {
						"s3:x-amz-acl": "bucket-owner-full-control"
						}
					},
					"Effect": "Allow",
					"Principal": {
						"Service": "cloudtrail.amazonaws.com"
					},
					"Resource": "arn:aws:s3:::turbot-111122224444-ap-northeast-1/AWSLogs/111122224444/*",
					"Sid": "AWSCloudTrailWrite"
					},
					{
					"Action": "s3:PutObject",
					"Effect": "Allow",
					"Principal": {
						"AWS": "arn:aws:iam::123456560864:root"
					},
					"Resource": "arn:aws:s3:::turbot-111122224444-ap-northeast-1/*",
					"Sid": "AWSELBWrite"
					},
					{
					"Action": "s3:GetBucketAcl",
					"Effect": "Allow",
					"Principal": {
						"AWS": "arn:aws:iam::123456285394:user/logs"
					},
					"Resource": "arn:aws:s3:::turbot-111122224444-ap-northeast-1",
					"Sid": "AWSRedshiftAclCheck"
					},
					{
					"Action": "s3:PutObject",
					"Effect": "Allow",
					"Principal": {
						"AWS": "arn:aws:iam::123456285394:user/logs"
					},
					"Resource": "arn:aws:s3:::turbot-111122224444-ap-northeast-1/*",
					"Sid": "AWSRedshiftWrite"
					},
					{
					"Action": "s3:GetBucketAcl",
					"Effect": "Allow",
					"Principal": {
						"Service": "delivery.logs.amazonaws.com"
					},
					"Resource": "arn:aws:s3:::turbot-111122224444-ap-northeast-1",
					"Sid": "AWSLogDeliveryAclCheck"
					},
					{
					"Action": "s3:PutObject",
					"Condition": {
						"StringEquals": {
						"s3:x-amz-acl": "bucket-owner-full-control"
						}
					},
					"Effect": "Allow",
					"Principal": {
						"Service": "delivery.logs.amazonaws.com"
					},
					"Resource": "arn:aws:s3:::turbot-111122224444-ap-northeast-1/AWSLogs/111122224444/*",
					"Sid": "AWSLogDeliveryWrite"
					}
				],
				"Version": "2012-10-17"
			}`,
			`{
        	"access_level": "public",
        	"allowed_organization_ids": null,
        	"allowed_principals": [
        		"arn:aws:iam::123456285394:user/logs",
        		"arn:aws:iam::123456560864:root",
        		"cloudtrail.amazonaws.com",
        		"delivery.logs.amazonaws.com"
        	],
        	"allowed_principal_account_ids": [
        		"123456285394",
        		"123456560864"
        	],
        	"allowed_principal_federated_identities": null,
        	"allowed_principal_services": [
        		"cloudtrail.amazonaws.com",
        		"delivery.logs.amazonaws.com"
        	],
        	"is_public": true,
        	"public_access_levels": [
        		"Read",
        		"Write"
        	],
        	"public_statement_ids": [
        		"AWSCloudTrailAclCheck",
        		"AWSCloudTrailWrite",
        		"AWSLogDeliveryAclCheck",
        		"AWSLogDeliveryWrite"
        	]
        }`,
		},
	}

	for index, test := range testCases {
		t.Run(strconv.Itoa(index+1), func(t *testing.T) {
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
func TestS3ExampleResourcePoliciesWithConditions(t *testing.T) {
	testCases := []publicAccessTest{
		{
			`Granting permissions to multiple accounts with added conditions`,
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
				"allowed_organization_ids": null,
				"allowed_principals": [
					"arn:aws:iam::111122223333:root",
					"arn:aws:iam::444455556666:root"
				],
				"allowed_principal_account_ids": [
					"444455556666"
				],
				"allowed_principal_federated_identities": null,
				"allowed_principal_services": null,
				"is_public": false,
				"public_access_levels": null,
				"public_statement_ids": null
			}`,
		},
		{
			`Limiting access to specific IP addresses`,
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
        	"allowed_organization_ids": null,
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
			`Policy allows request with mfa only`,
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
        	"public_access_levels": ["List","Read"],
        	"public_statement_ids": [
        		"Statement1"
        	]
        }`,
		},
		{
			`Policy allows request with mfa and principal account limitation`,
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
			`Granting permissions for Amazon S3 Inventory and Amazon S3 analytics with only source arn condition`,
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
        	"access_level": "public",
        	"allowed_organization_ids": null,
        	"allowed_principals": [
        		"arn:aws:s3:::sourcebucket",
        		"s3.amazonaws.com"
        	],
        	"allowed_principal_account_ids": null,
        	"allowed_principal_federated_identities": null,
        	"allowed_principal_services": [
        		"s3.amazonaws.com"
        	],
        	"is_public": true,
        	"public_access_levels": [
        		"Write"
        	],
        	"public_statement_ids": [
        		"InventoryAndAnalyticsExamplePolicy"
        	]
        }`,
		},
	}

	for index, test := range testCases {
		t.Run(strconv.Itoa(index+1), func(t *testing.T) {
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

func TestPublicPolicies(t *testing.T) {
	testCases := []string{
		`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Sid": "Allow Member Account Access",
					"Effect": "Allow",
					"Principal": "*",
					"Action": "s3:GetObject",
					"Resource": "arn:aws:s3:::sample-bucket/*"
				}
			]
		}`,
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
					"Resource": "arn:aws:sns:us-east-1:111122225555:smyth-test-2"
				}
			]
		}`,
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
			"Version": "2012-10-17",
			"Statement": [
				{
					"Sid": "SharedAccess4",
					"Effect": "Allow",
					"Principal": "*",
					"Action": [
						"s3:listbucket",
						"s3:getobject"
					],
					"Resource": [
						"arn:aws:s3:::test-anonymous-access",
						"arn:aws:s3:::test-anonymous-access/*"
					],
					"Condition": {
						"StringLike": {
							"aws:PrincipalOrgID": ["o-1a2b3c4d*"]
						}
					}
				}
			]
		}`,
		`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Sid": "ServicePublicAccess",
					"Effect": "Allow",
					"Principal": {
						"Service": "cloudtrail.amazonaws.com"
					},
					"Action": "SNS:Publish",
					"Resource": "arn:aws:sns:region:SNSTopicOwnerAccountId:SNSTopicName"
				}
			]
		}`,
		`{
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {
						"Service": "s3.amazonaws.com"
					},
					"Action": "sns:Publish",
					"Resource": "arn:aws:sns:us-east-2:111122225555:MyTopic",
					"Condition": {
						"ArnLike": {
							"aws:SourceArn": "arn:aws:S3:::test-bucket"
						}
					}
				}
			]
		}`,
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
			"Version": "2012-10-17",
			"Statement": [
				{
					"Sid": "SharedAccess4",
					"Effect": "Allow",
					"Principal": {
						"AWS": "*"
					},
					"Action": [
						"s3:listbucket",
						"s3:getobject"
					],
					"Resource": [
						"arn:aws:s3:::test-anonymous-access",
						"arn:aws:s3:::test-anonymous-access/*"
					],
					"Condition": {
						"StringLike": {
							"aws:PrincipalAccount": [
								"11112222333*",
								"11112222444*"
							]
						}
					}
				}
			]
		}`,
	}

	for index, policy := range testCases {
		t.Run(fmt.Sprintf("%d", index+1), func(t *testing.T) {
			policy, err := canonicalPolicy(policy)
			if err != nil {
				t.Errorf("Test: %d Policy canonicalization failed with error: %#v\n", index+1, err)
			}

			policyObject, ok := policy.(Policy)
			if !ok {
				t.Errorf("Test: %d Policy coercion failed with error: %#v\n", index+1, err)
			}
			evaluatedObj, err := policyObject.EvaluatePolicy("111122225555")
			if err != nil {
				t.Errorf("Test: %d\nPolicy evaluation failed with error: %#v\n", index+1, err)
			}

			if !evaluatedObj.IsPublic {
				strdata, _ := json.MarshalIndent(evaluatedObj, "", "\t")
				t.Errorf("Expected %d to be public but it is %v\n", index+1, string(strdata))
			}
		})
	}
}

func TestSharedPolicies(t *testing.T) {
	testCases := []string{
		`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Sid": "SharedAccess1",
					"Effect": "Allow",
					"Principal": {
						"AWS": [
							"arn:aws:iam::111122223333:root",
							"arn:aws:iam::111122224444:root"
						]
					},
					"Action": "s3:ListBucket",
					"Resource": "arn:aws:s3:::test-anonymous-access"
				}
			]
		}`,
		`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Sid": "SharedAccess2",
					"Effect": "Allow",
					"Principal": "*",
					"Action": "s3:ListBucket",
					"Resource": "arn:aws:s3:::test-anonymous-access",
					"Condition": {
						"StringEquals": {
							"aws:PrincipalAccount": ["111122221111", "111122223333"]
						}
					}
				}
			]
		}`,
		`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Sid": "SharedAccess3",
					"Effect": "Allow",
					"Principal": "*",
					"Action": ["s3:listbucket", "s3:GetObject"],
					"Resource": [
						"arn:aws:s3:::test-anonymous-access/*",
						"arn:aws:s3:::test-anonymous-access"
					],
					"Condition": {
						"StringEquals": {
							"aws:PrincipalArn": "arn:aws:iam::013122550996:root"
						}
					}
				}
			]
		}`,
		`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Sid": "SharedAccess4",
					"Effect": "Allow",
					"Principal": "*",
					"Action": [
						"s3:listbucket",
						"s3:getobject"
					],
					"Resource": [
						"arn:aws:s3:::test-anonymous-access",
						"arn:aws:s3:::test-anonymous-access/*"
					],
					"Condition": {
						"StringEquals": {
							"aws:PrincipalOrgID": ["o-1a2b3c4d5e"]
						}
					}
				}
			]
		}`,
		`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Sid": "SharedAccess4",
					"Effect": "Allow",
					"Principal": "*",
					"Action": [
						"s3:listbucket",
						"s3:getobject"
					],
					"Resource": [
						"arn:aws:s3:::test-anonymous-access",
						"arn:aws:s3:::test-anonymous-access/*"
					],
					"Condition": {
						"StringLike": {
							"aws:PrincipalOrgID": ["o-1a2b3c4d5e"]
						}
					}
				}
			]
		}`,
		`{
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {
						"Service": "someservice.amazonaws.com"
					},
					"Action": "sns:Publish",
					"Resource": "arn:aws:sns:us-east-2:111122225555:MyTopic",
					"Condition": {
						"StringEquals": {
							"aws:SourceAccount": "444455556666"
						}
					}
				}
			]
		}`,
		// https://docs.aws.amazon.com/IAM/latest/UserGuide/confused-deputy.html#cross-service-confused-deputy-prevention
		`{
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {
						"Service": "s3.amazonaws.com"
					},
					"Action": "sns:Publish",
					"Resource": "arn:aws:sns:us-east-2:111122225555:MyTopic",
					"Condition": {
						"StringEquals": {
							"aws:SourceAccount": "444455556666"
						},
						"ArnLike": {
							"aws:SourceArn": "arn:aws:S3:::test-bucket"
						}
					}
				}
			]
		}`,
		`{
			"Version": "2012-10-17",
			"Statement": {
				"Effect": "Allow",
				"Principal": {
					"Service": "ssm-incidents.amazonaws.com"
				},
				"Action": "sts:AssumeRole",
				"Condition": {
					"ArnLike": {
						"aws:SourceArn": "arn:aws:ssm-incidents:*:111122223333:incident-record/myresponseplan/*"
					}
				}
			}
		}`,
	}

	for index, policy := range testCases {
		t.Run(fmt.Sprintf("%d", index+1), func(t *testing.T) {
			policy, err := canonicalPolicy(policy)
			if err != nil {
				t.Errorf("Test: %d Policy canonicalization failed with error: %#v\n", index+1, err)
			}

			policyObject, ok := policy.(Policy)
			if !ok {
				t.Errorf("Test: %d Policy coercion failed with error: %#v\n", index+1, err)
			}
			evaluatedObj, err := policyObject.EvaluatePolicy("111122225555")
			if err != nil {
				t.Errorf("Test: %d\nPolicy evaluation failed with error: %#v\n", index+1, err)
			}

			strdata, _ := json.MarshalIndent(evaluatedObj, "", "\t")
			if evaluatedObj.AccessLevel != "shared" {
				t.Errorf("Expected %d to be shared but it is %s. Ealuation: \n%v\n", index+1, evaluatedObj.AccessLevel, string(strdata))
			}
			// else {
			// 	fmt.Println("Evaluation: \n", string(strdata))
			// }
		})
	}
}

func TestPrivatePolicies(t *testing.T) {
	testCases := []string{
		`{
				"Statement": [
					{
						"Effect": "Allow",
						"Principal": "*",
						"Action": "SNS:Publish",
						"Resource": "arn:aws:sns:us-east-2:123456789012:MyTopic",
						"Condition": {
							"ArnLike": {
								"aws:SourceArn": "arn:aws:cloudwatch:us-east-2:123456789012:alarm:*"
							}
						}
					}
				]
			}`,
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
							"ArnLike": {
	             	"aws:SourceArn": "arn:aws:iam::123456789012:users/*"
	           	}
						}
					}
				]
			}`,
		`{
				"Version":"2012-10-17",
				"Statement":[
					{
						"Sid":"PolicyForAllowUploadWithACL",
						"Effect":"Allow",
						"Principal":{"AWS":"123456789012"},
						"Action":"s3:PutObject",
						"Resource":"arn:aws:s3:::DOC-EXAMPLE-BUCKET/*",
						"Condition": {
							"StringEquals": {"s3:x-amz-acl":"bucket-owner-full-control"}
						}
					}
				]
			}`,
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
								"aws:SourceAccount": "123456789012",
								"s3:x-amz-acl": "bucket-owner-full-control"
							}
						}
					}
				]
			}`,
	}

	for index, policy := range testCases {
		t.Run(fmt.Sprintf("%d", index+1), func(t *testing.T) {
			policy, err := canonicalPolicy(policy)
			if err != nil {
				t.Errorf("Test: %d Policy canonicalization failed with error: %#v\n", index+1, err)
			}

			policyObject, ok := policy.(Policy)
			if !ok {
				t.Errorf("Test: %d Policy coercion failed with error: %#v\n", index+1, err)
			}
			evaluatedObj, err := policyObject.EvaluatePolicy("123456789012")
			if err != nil {
				t.Errorf("Test: %d\nPolicy evaluation failed with error: %#v\n", index+1, err)
			}

			if evaluatedObj.AccessLevel != "private" {
				strdata, _ := json.MarshalIndent(evaluatedObj, "", "\t")
				t.Errorf("Expected %d to be private but it is %v\n", index+1, string(strdata))
			}
		})
	}
}

func TestPoliciesForArnOperators(t *testing.T) {
	testCases := []struct {
		Policy      string
		AccessLevel string
	}{
		{
			`{
			"Version": "2012-10-17",
			"Id": "default",
			"Statement": [
				{
					"Sid": "s3-event-cezary_for_lambda-s3",
					"Effect": "Allow",
					"Principal": {
						"Service": "s3.amazonaws.com"
					},
					"Action": "lambda:InvokeFunction",
					"Resource": "arn:aws:lambda:eu-west-1:123456789012:function:lambda-s3",
					"Condition": {
						"StringEquals": {
							"AWS:SourceAccount": "123456789012"
						},
						"ArnLike": {
							"AWS:SourceArn": "arn:aws:s3:::test-bucket-cezary"
						}
					}
				}
			]
		}`,
			"private",
		},
		{
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
							"ArnLike": {
								"aws:PrincipalArn": "arn:aws:iam::*:user/*"
							}
						}
					}
				]
			}`,
			"public",
		},
		{
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
			"shared",
		},
	}

	for index, item := range testCases {
		t.Run(fmt.Sprintf("%d", index+1), func(t *testing.T) {
			policy, err := canonicalPolicy(item.Policy)
			if err != nil {
				t.Errorf("Test: %d Policy canonicalization failed with error: %#v\n", index+1, err)
			}

			policyObject, ok := policy.(Policy)
			if !ok {
				t.Errorf("Test: %d Policy coercion failed with error: %#v\n", index+1, err)
			}
			evaluatedObj, err := policyObject.EvaluatePolicy("123456789012")
			if err != nil {
				t.Errorf("Test: %d\nPolicy evaluation failed with error: %#v\n", index+1, err)
			}

			if evaluatedObj.AccessLevel != item.AccessLevel {
				strdata, _ := json.MarshalIndent(evaluatedObj, "", "\t")
				fmt.Println("policy:", string(strdata))
				t.Errorf("Expected %d to be %s but it is %s\n", index+1, item.AccessLevel, evaluatedObj.AccessLevel)
			}
		})
	}
}

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
			// fmt.Println("start", i+1, time.Now())
			accessLevels := GetAccessLevelsFromActions(permissionsData, test.Actions)
			// fmt.Println("end", i+1, time.Now())

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
