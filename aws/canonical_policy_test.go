package aws

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestConvertPolicySortAndDups(t *testing.T) {
	testCase := string(
		//// sort things, remove duplicates
		`{
			"Version": "2012-10-17",
			"Id": "TestSortAndRemoveDuplicates",
			"Statement": [
				{
					"Effect": "Allow",
					"Action": [
						"C",
						"A",
						"B",
						"b",
						"c",
						"a",
						"2",
						"1"
					],
					"NotAction": [
						"C",
						"A",
						"B",
						"b",
						"c",
						"a",
						"2",
						"1"
					],
					"Principal": {
						"AWS": [
							"C",
							"A",
							"B"
						],
						"Service": [
							"C",
							"A",
							"B"
						],
						"Federated": [
							"C",
							"A",
							"B"
						]
					},
					"NotPrincipal": {
						"Federated": [
							"C",
							"A",
							"B"
						],
						"AWS": [
							"C",
							"A",
							"B"
						],
						"Service": [
							"C",
							"A",
							"B"
						]
					},
					"Resource": [
						"C",
						"A",
						"B",
						"b",
						"c",
						"a",
						"2",
						"1"
					],
					"NotResource": [
						"C",
						"A",
						"B",
						"b",
						"c",
						"a",
						"2",
						"1"
					],
					"Condition": {
						"z": {
							"a": false
						},
						"y": {
							"a": [
								"b",
								"c"
							]
						},
						"x": {
							"b": "b"
						},
						"a": {
							"b": "b",
							"c": "c",
							"d": "d",
							"e": "e",
							"a": [
								"c",
								"b",
								"a"
							],
							"g": "g",
							"f": "f"

						}

					}
				}
			]
		}`)

	pol, err := canonicalPolicy(testCase)
	if err != nil {
		t.Errorf("Convert failed for case '%s': %v", pol, err)
	}
	prettyPrint(pol)
}

func TestConvertPolicyWithSingleStatement(t *testing.T) {
	testCase := string(
		// single statement policies
		`{
			"Version": "2012-10-17",
			"Statement": {
				"Effect": "Allow",
				"Action": [
					"acm:DescribeCertificate",
					"acm:ListCertificates",
					"acm:GetCertificate",
					"acm:ListTagsForCertificate"
				],
				"Resource": "*"
			}
		}`)

	pol, err := canonicalPolicy(testCase)
	if err != nil {
		t.Errorf("Convert failed for case '%s': %v", pol, err)
	}
	prettyPrint(pol)
}

func TestConvertPolicyWithBoolsAndInts(t *testing.T) {
	testCase := string(
		// s/// boolean, int in condition
		`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Action": [
						"ec2:ModifyVpcEndpoint",
						"ec2:DeleteVpcEndpoints"
					],
					"Resource": "arn:aws:ec2:*:*:vpc-endpoint/*",
					"Condition": {
						"Null": {
							"aws:ResourceTag/AmazonMWAAManaged": false
						},
						"StringEquals": {
							"ec2:CreateAction": [
								"CreateVpcEndpoint",
								"DeleteVpcEndpoint"
							]
						},
						"ForAnyValue:StringEquals": {
							"aws:TagKeys": "AmazonMWAAManaged"
						},
						"ThisIsJustToTestConversion": {
							"int_as_string": "42",
							"int_as_int": 42,
							"bool_as_string_true": "true",
							"bool_as_bool_true": true,
							"bool_as_string_false": "false",
							"bool_as_bool_false": false,
							"string": "this is a string"

						}

					}
				}
			]
		}`)

	pol, err := canonicalPolicy(testCase)
	if err != nil {
		t.Errorf("Convert failed for case '%s': %v", pol, err)
	}
	prettyPrint(pol)
}

func TestConvertPolicy(t *testing.T) {
	cases := []string{

		// Regressions
		`{ "Version": "2012-10-17", "Statement": [ { "Effect": "Allow", "Action": [ "logs:CreateLogStream", "logs:CreateLogGroup", "logs:DescribeLogGroups" ], "Resource": "arn:aws:logs:*:*:log-group:airflow-*:*" }, { "Effect": "Allow", "Action": [ "ec2:AttachNetworkInterface", "ec2:CreateNetworkInterface", "ec2:CreateNetworkInterfacePermission", "ec2:DeleteNetworkInterface", "ec2:DeleteNetworkInterfacePermission", "ec2:DescribeDhcpOptions", "ec2:DescribeNetworkInterfaces", "ec2:DescribeSecurityGroups", "ec2:DescribeSubnets", "ec2:DescribeVpcEndpoints", "ec2:DescribeVpcs", "ec2:DetachNetworkInterface" ], "Resource": "*" }, { "Effect": "Allow", "Action": "ec2:CreateVpcEndpoint", "Resource": "arn:aws:ec2:*:*:vpc-endpoint/*", "Condition": { "ForAnyValue:StringEquals": { "aws:TagKeys": "AmazonMWAAManaged" } } }, { "Effect": "Allow", "Action": [ "ec2:ModifyVpcEndpoint", "ec2:DeleteVpcEndpoints" ], "Resource": "arn:aws:ec2:*:*:vpc-endpoint/*", "Condition": { "Null": { "aws:ResourceTag/AmazonMWAAManaged": false } } }, { "Effect": "Allow", "Action": [ "ec2:CreateVpcEndpoint", "ec2:ModifyVpcEndpoint" ], "Resource": [ "arn:aws:ec2:*:*:vpc/*", "arn:aws:ec2:*:*:security-group/*", "arn:aws:ec2:*:*:subnet/*" ] }, { "Effect": "Allow", "Action": "ec2:CreateTags", "Resource": "arn:aws:ec2:*:*:vpc-endpoint/*", "Condition": { "StringEquals": { "ec2:CreateAction": "CreateVpcEndpoint" }, "ForAnyValue:StringEquals": { "aws:TagKeys": "AmazonMWAAManaged" } } } ] }`,

		// //// S3 bucket policies
		`{"Version":"2012-10-17","Statement":[{"Sid":"AWSCloudTrailAclCheck","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"s3:GetBucketAcl","Resource":"arn:aws:s3:::turbot-876515858155-sa-east-1"},{"Sid":"AWSCloudTrailWrite","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"s3:PutObject","Resource":"arn:aws:s3:::turbot-876515858155-sa-east-1/AWSLogs/876515858155/*","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}}},{"Sid":"AWSELBWrite","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::507241528517:root"},"Action":"s3:PutObject","Resource":"arn:aws:s3:::turbot-876515858155-sa-east-1/*"},{"Sid":"AWSRedshiftAclCheck","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::075028567923:user/logs"},"Action":"s3:GetBucketAcl","Resource":"arn:aws:s3:::turbot-876515858155-sa-east-1"},{"Sid":"AWSRedshiftWrite","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::075028567923:user/logs"},"Action":"s3:PutObject","Resource":"arn:aws:s3:::turbot-876515858155-sa-east-1/*"},{"Sid":"AWSLogDeliveryAclCheck","Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Action":"s3:GetBucketAcl","Resource":"arn:aws:s3:::turbot-876515858155-sa-east-1"},{"Sid":"AWSLogDeliveryWrite","Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Action":"s3:PutObject","Resource":"arn:aws:s3:::turbot-876515858155-sa-east-1/AWSLogs/876515858155/*","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}}}]}`,

		`{"Id":"Policy1590074992386","Statement":[{"Action":"s3:listBucket","Condition":{"StringEquals":{"AWS:anything":["ACCOUNT-ID"]}},"Effect":"Allow","Principal":{"AWS":["AROASNJZ76JBEFDBLBIMV","arn:aws:iam::039305405804:root","arn:aws:iam::166014743106:user/jsmyth","arn:aws:iam::235268162285:root"],"Federated":["literally anything","accounts.google.com","arn:aws:iam::AWS-account-ID:saml-provider/provider-name","graph.facebook.com","cognito-identity.amazonaws.com"],"Service":["cloudtrail.amazonaws.com","sns.amazonaws.com"]},"Resource":"arn:aws:s3:::JSMYTH-test-bucket-8765","Sid":"Stmt1590074983320"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":"s3:LISTBUCKET","Effect":"Allow","Principal":"*","Resource":"arn:aws:s3:::vandelay-INSECURE-test-bucket-do-not-use","Sid":"Stmt1600291154570"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":"s3:*","Condition":{"Bool":{"aws:SecureTransport":"false"}},"Effect":"Deny","Principal":"*","Resource":["arn:aws:s3:::vandelay-industries-elaines-bucket","arn:aws:s3:::vandelay-industries-elaines-bucket/*"],"Sid":"MustBeEncryptedInTransit"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-us-east-1","Sid":"AWSCloudTrailAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-us-east-1/AWSLogs/876515858155/*","Sid":"AWSCloudTrailWrite"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::127311923021:root"},"Resource":"arn:aws:s3:::turbot-876515858155-us-east-1/*","Sid":"AWSELBWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::193672423079:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-us-east-1","Sid":"AWSRedshiftAclCheck"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::193672423079:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-us-east-1/*","Sid":"AWSRedshiftWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-us-east-1","Sid":"AWSLogDeliveryAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-us-east-1/AWSLogs/876515858155/*","Sid":"AWSLogDeliveryWrite"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-us-east-2","Sid":"AWSCloudTrailAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-us-east-2/AWSLogs/876515858155/*","Sid":"AWSCloudTrailWrite"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::033677994240:root"},"Resource":"arn:aws:s3:::turbot-876515858155-us-east-2/*","Sid":"AWSELBWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::391106570357:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-us-east-2","Sid":"AWSRedshiftAclCheck"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::391106570357:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-us-east-2/*","Sid":"AWSRedshiftWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-us-east-2","Sid":"AWSLogDeliveryAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-us-east-2/AWSLogs/876515858155/*","Sid":"AWSLogDeliveryWrite"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-ca-central-1","Sid":"AWSCloudTrailAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-ca-central-1/AWSLogs/876515858155/*","Sid":"AWSCloudTrailWrite"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::985666609251:root"},"Resource":"arn:aws:s3:::turbot-876515858155-ca-central-1/*","Sid":"AWSELBWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::907379612154:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-ca-central-1","Sid":"AWSRedshiftAclCheck"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::907379612154:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-ca-central-1/*","Sid":"AWSRedshiftWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-ca-central-1","Sid":"AWSLogDeliveryAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-ca-central-1/AWSLogs/876515858155/*","Sid":"AWSLogDeliveryWrite"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":"s3:*","Condition":{"Bool":{"aws:SecureTransport":"false"}},"Effect":"Deny","Principal":"*","Resource":["arn:aws:s3:::vandelay-industries-darins-bucket","arn:aws:s3:::vandelay-industries-darins-bucket/*"],"Sid":"MustBeEncryptedInTransit"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-us-west-2","Sid":"AWSCloudTrailAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-us-west-2/AWSLogs/876515858155/*","Sid":"AWSCloudTrailWrite"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::797873946194:root"},"Resource":"arn:aws:s3:::turbot-876515858155-us-west-2/*","Sid":"AWSELBWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::902366379725:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-us-west-2","Sid":"AWSRedshiftAclCheck"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::902366379725:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-us-west-2/*","Sid":"AWSRedshiftWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-us-west-2","Sid":"AWSLogDeliveryAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-us-west-2/AWSLogs/876515858155/*","Sid":"AWSLogDeliveryWrite"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-west-2","Sid":"AWSCloudTrailAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-west-2/AWSLogs/876515858155/*","Sid":"AWSCloudTrailWrite"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::652711504416:root"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-west-2/*","Sid":"AWSELBWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::307160386991:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-west-2","Sid":"AWSRedshiftAclCheck"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::307160386991:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-west-2/*","Sid":"AWSRedshiftWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-west-2","Sid":"AWSLogDeliveryAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-west-2/AWSLogs/876515858155/*","Sid":"AWSLogDeliveryWrite"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-west-3","Sid":"AWSCloudTrailAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-west-3/AWSLogs/876515858155/*","Sid":"AWSCloudTrailWrite"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::009996457667:root"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-west-3/*","Sid":"AWSELBWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::915173422425:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-west-3","Sid":"AWSRedshiftAclCheck"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::915173422425:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-west-3/*","Sid":"AWSRedshiftWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-west-3","Sid":"AWSLogDeliveryAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-west-3/AWSLogs/876515858155/*","Sid":"AWSLogDeliveryWrite"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-west-1","Sid":"AWSCloudTrailAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-west-1/AWSLogs/876515858155/*","Sid":"AWSCloudTrailWrite"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::156460612806:root"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-west-1/*","Sid":"AWSELBWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::210876761215:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-west-1","Sid":"AWSRedshiftAclCheck"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::210876761215:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-west-1/*","Sid":"AWSRedshiftWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-west-1","Sid":"AWSLogDeliveryAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-west-1/AWSLogs/876515858155/*","Sid":"AWSLogDeliveryWrite"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-central-1","Sid":"AWSCloudTrailAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-central-1/AWSLogs/876515858155/*","Sid":"AWSCloudTrailWrite"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::054676820928:root"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-central-1/*","Sid":"AWSELBWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::053454850223:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-central-1","Sid":"AWSRedshiftAclCheck"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::053454850223:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-central-1/*","Sid":"AWSRedshiftWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-central-1","Sid":"AWSLogDeliveryAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-central-1/AWSLogs/876515858155/*","Sid":"AWSLogDeliveryWrite"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-northeast-1","Sid":"AWSCloudTrailAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-northeast-1/AWSLogs/876515858155/*","Sid":"AWSCloudTrailWrite"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::582318560864:root"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-northeast-1/*","Sid":"AWSELBWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::404641285394:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-northeast-1","Sid":"AWSRedshiftAclCheck"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::404641285394:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-northeast-1/*","Sid":"AWSRedshiftWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-northeast-1","Sid":"AWSLogDeliveryAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-northeast-1/AWSLogs/876515858155/*","Sid":"AWSLogDeliveryWrite"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-north-1","Sid":"AWSCloudTrailAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-north-1/AWSLogs/876515858155/*","Sid":"AWSCloudTrailWrite"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::897822967062:root"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-north-1/*","Sid":"AWSELBWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::729911121831:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-north-1","Sid":"AWSRedshiftAclCheck"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::729911121831:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-north-1/*","Sid":"AWSRedshiftWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-north-1","Sid":"AWSLogDeliveryAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-eu-north-1/AWSLogs/876515858155/*","Sid":"AWSLogDeliveryWrite"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-sa-east-1","Sid":"AWSCloudTrailAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-sa-east-1/AWSLogs/876515858155/*","Sid":"AWSCloudTrailWrite"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::507241528517:root"},"Resource":"arn:aws:s3:::turbot-876515858155-sa-east-1/*","Sid":"AWSELBWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::075028567923:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-sa-east-1","Sid":"AWSRedshiftAclCheck"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::075028567923:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-sa-east-1/*","Sid":"AWSRedshiftWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-sa-east-1","Sid":"AWSLogDeliveryAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-sa-east-1/AWSLogs/876515858155/*","Sid":"AWSLogDeliveryWrite"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-northeast-2","Sid":"AWSCloudTrailAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-northeast-2/AWSLogs/876515858155/*","Sid":"AWSCloudTrailWrite"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::600734575887:root"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-northeast-2/*","Sid":"AWSELBWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::760740231472:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-northeast-2","Sid":"AWSRedshiftAclCheck"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::760740231472:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-northeast-2/*","Sid":"AWSRedshiftWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-northeast-2","Sid":"AWSLogDeliveryAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-northeast-2/AWSLogs/876515858155/*","Sid":"AWSLogDeliveryWrite"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-southeast-2","Sid":"AWSCloudTrailAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-southeast-2/AWSLogs/876515858155/*","Sid":"AWSCloudTrailWrite"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::783225319266:root"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-southeast-2/*","Sid":"AWSELBWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::762762565011:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-southeast-2","Sid":"AWSRedshiftAclCheck"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::762762565011:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-southeast-2/*","Sid":"AWSRedshiftWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-southeast-2","Sid":"AWSLogDeliveryAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-southeast-2/AWSLogs/876515858155/*","Sid":"AWSLogDeliveryWrite"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-south-1","Sid":"AWSCloudTrailAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-south-1/AWSLogs/876515858155/*","Sid":"AWSCloudTrailWrite"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::718504428378:root"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-south-1/*","Sid":"AWSELBWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::865932855811:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-south-1","Sid":"AWSRedshiftAclCheck"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::865932855811:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-south-1/*","Sid":"AWSRedshiftWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-south-1","Sid":"AWSLogDeliveryAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-south-1/AWSLogs/876515858155/*","Sid":"AWSLogDeliveryWrite"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-southeast-1","Sid":"AWSCloudTrailAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-southeast-1/AWSLogs/876515858155/*","Sid":"AWSCloudTrailWrite"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::114774131450:root"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-southeast-1/*","Sid":"AWSELBWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::361669875840:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-southeast-1","Sid":"AWSRedshiftAclCheck"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::361669875840:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-southeast-1/*","Sid":"AWSRedshiftWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-southeast-1","Sid":"AWSLogDeliveryAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-ap-southeast-1/AWSLogs/876515858155/*","Sid":"AWSLogDeliveryWrite"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-us-west-1","Sid":"AWSCloudTrailAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-us-west-1/AWSLogs/876515858155/*","Sid":"AWSCloudTrailWrite"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::027434742980:root"},"Resource":"arn:aws:s3:::turbot-876515858155-us-west-1/*","Sid":"AWSELBWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::262260360010:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-us-west-1","Sid":"AWSRedshiftAclCheck"},{"Action":"s3:PutObject","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::262260360010:user/logs"},"Resource":"arn:aws:s3:::turbot-876515858155-us-west-1/*","Sid":"AWSRedshiftWrite"},{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-us-west-1","Sid":"AWSLogDeliveryAclCheck"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"delivery.logs.amazonaws.com"},"Resource":"arn:aws:s3:::turbot-876515858155-us-west-1/AWSLogs/876515858155/*","Sid":"AWSLogDeliveryWrite"}],"Version":"2012-10-17"}`,

		/// IAM policies
		`{"Statement":[{"Action":["acm:*"],"Effect":"Allow","Resource":"*","Sid":"AWSacmAdmin"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":["dynamodb:*"],"Effect":"Allow","Resource":"*","Sid":"AWSdynamodbOwner"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":["acm:listc*","application-autoscaling:der*","application-autoscaling:des*","application-autoscaling:r*","autoscaling:a*","autoscaling:co*","autoscaling:createorupdatet*","autoscaling:deleten*","autoscaling:deleteta*","autoscaling:des*","autoscaling:det*","autoscaling:di*","autoscaling:e*","autoscaling:putn*","autoscaling:r*","autoscaling:seti*","autoscaling:su*","autoscaling:t*","ec2:copys*","ec2:createsn*","ec2:createta*","ec2:deleteta*","ec2:describea*","ec2:describeb*","ec2:describecl*","ec2:describeco*","ec2:describeel*","ec2:describeex*","ec2:describefp*","ec2:describeh*","ec2:describeia*","ec2:describeid*","ec2:describeim*","ec2:describeins*","ec2:describek*","ec2:describel*","ec2:describem*","ec2:describenetworki*","ec2:describepl*","ec2:describere*","ec2:describesc*","ec2:describese*","ec2:describesn*","ec2:describesp*","ec2:describest*","ec2:describet*","ec2:describevo*","ec2:getc*","ec2:geth*","ec2:getl*","ec2:getr*","ec2:reb*","ec2:repo*","ec2:st*","elasticloadbalancing:addt*","elasticloadbalancing:der*","elasticloadbalancing:des*","elasticloadbalancing:reg*","elasticloadbalancing:removet*"],"Effect":"Allow","Resource":"*","Sid":"AWSec2Operator"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":["logs:*"],"Effect":"Allow","Resource":"*","Sid":"AWSlogsAdmin"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":["dynamodb:b*","dynamodb:c*","dynamodb:d*","dynamodb:g*","dynamodb:l*","dynamodb:put*","dynamodb:q*","dynamodb:r*","dynamodb:s*","dynamodb:t*","dynamodb:u*"],"Effect":"Allow","Resource":"*","Sid":"AWSdynamodbAdmin"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":["ec2:describeac*","ec2:describeav*","ec2:describesecuritygroups*","ec2:describevpcs*","elasticache:ad*","elasticache:co*","elasticache:createcachec*","elasticache:createcachep*","elasticache:createcachesu*","elasticache:creater*","elasticache:creates*","elasticache:dec*","elasticache:deletecachec*","elasticache:deletecachep*","elasticache:deletecachesu*","elasticache:deleter*","elasticache:deletes*","elasticache:des*","elasticache:i*","elasticache:l*","elasticache:m*","elasticache:reb*","elasticache:rem*","elasticache:res*","elasticache:t*","sns:listsubscriptions","sns:listto*"],"Effect":"Allow","Resource":"*","Sid":"AWSelastiCacheAdmin"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":["elasticfilesystem:*"],"Effect":"Allow","Resource":"*","Sid":"AWSefsAdmin"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":["acm-pca:*","acm:*","amplify:*","apigateway:*","application-autoscaling:*","appmesh:*","athena:*","autoscaling:a*","autoscaling:co*","autoscaling:createorupdatet*","autoscaling:deleten*","autoscaling:deleteta*","autoscaling:des*","autoscaling:det*","autoscaling:di*","autoscaling:e*","autoscaling:putn*","autoscaling:r*","autoscaling:seti*","autoscaling:su*","autoscaling:t*","aws-marketplace-management:v*","aws-marketplace:b*","aws-marketplace:g*","aws-marketplace:m*","aws-marketplace:r*","aws-marketplace:v*","aws-portal:v*","backup-storage:*","backup:*","ce:*","cloudformation:ca*","cloudformation:co*","cloudformation:createc*","cloudformation:createstack","cloudformation:createu*","cloudformation:deletec*","cloudformation:deletestack","cloudformation:des*","cloudformation:det*","cloudformation:e*","cloudformation:g*","cloudformation:l*","cloudformation:se*","cloudformation:si*","cloudformation:updatestack","cloudformation:updatet*","cloudfront:*","cloudsearch:*","cloudtrail:*","cloudwatch:*","codebuild:*","codecommit:*","config:*","dax:*","ds:describedi*","dynamodb:b*","dynamodb:c*","dynamodb:d*","dynamodb:g*","dynamodb:l*","dynamodb:put*","dynamodb:q*","dynamodb:r*","dynamodb:s*","dynamodb:t*","dynamodb:u*","ec2-reports:*","ec2:acceptr*","ec2:al*","ec2:assi*","ec2:associatea*","ec2:associatei*","ec2:attachn*","ec2:attachvo*","ec2:cancels*","ec2:copys*","ec2:createk*","ec2:createl*","ec2:createnetworki*","ec2:createp*","ec2:createsn*","ec2:createsp*","ec2:createta*","ec2:createvo*","ec2:deletek*","ec2:deletel*","ec2:deletenetworki*","ec2:deletep*","ec2:deletesn*","ec2:deletesp*","ec2:deleteta*","ec2:deletevo*","ec2:describea*","ec2:describeb*","ec2:describec*","ec2:described*","ec2:describeel*","ec2:describeex*","ec2:describefp*","ec2:describeh*","ec2:describei*","ec2:describek*","ec2:describel*","ec2:describem*","ec2:describene*","ec2:describepl*","ec2:describer*","ec2:describes*","ec2:describet*","ec2:describevo*","ec2:describevpca*","ec2:describevpcpeeringconnection","ec2:describevpcs*","ec2:describevpn*","ec2:detachn*","ec2:detachvo*","ec2:disassociatea*","ec2:disassociatei*","ec2:enablevo*","ec2:g*","ec2:importk*","ec2:modifyh*","ec2:modifyid*","ec2:modifyin*","ec2:modifyl*","ec2:modifyn*","ec2:modifyr*","ec2:modifysn*","ec2:modifysp*","ec2:modifyvo*","ec2:mon*","ec2:purchaseh*","ec2:purchaser*","ec2:reb*","ec2:rel*","ec2:replacei*","ec2:repo*","ec2:req*","ec2:resetin*","ec2:resets*","ec2:ru*","ec2:st*","ec2:t*","ec2:un*","ec2messages:*","ecr:b*","ecr:c*","ecr:d*","ecr:g*","ecr:i*","ecr:l*","ecr:p*","ecr:st*","ecr:t*","ecr:u*","ecs:c*","ecs:de*","ecs:l*","ecs:pu*","ecs:registert*","ecs:ru*","ecs:startta*","ecs:sto*","ecs:t*","ecs:u*","eks:*","elasticache:ad*","elasticache:co*","elasticache:createcachec*","elasticache:createcachep*","elasticache:createcachesu*","elasticache:creater*","elasticache:creates*","elasticache:dec*","elasticache:deletecachec*","elasticache:deletecachep*","elasticache:deletecachesu*","elasticache:deleter*","elasticache:deletes*","elasticache:des*","elasticache:i*","elasticache:l*","elasticache:m*","elasticache:reb*","elasticache:rem*","elasticache:res*","elasticache:t*","elasticbeanstalk:*","elasticfilesystem:*","elasticloadbalancing:*","elasticmapreduce:*","es:a*","es:c*"],"Effect":"Allow","Resource":"*","Sid":"AWSAdmin1"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":["athena:*"],"Effect":"Allow","Resource":"*","Sid":"AWSathenaAdmin"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":["health:*","waf:*","route53:*","cloudFront:*","iam:*","sts:*"],"Effect":"Allow","Resource":"*"},{"Action":["xray:*","workspaces:*","snowball:*","servicecatalog:*","route53resolver:*","ram:*","glue:*","dms:*","datapipeline:*","comprehend:*","codedeploy:*","cloudhsm:*","batch:*","artifact:*","waf-regional:*","states:*","sqs:*","shield:*","securityHub:*","secretsManager:*","roboMaker:*","redshift:*","kafka:*","mq:*","lambda:*","kinesis:*","inspector:*","guardDuty:*","glacier:*","fsx:*","elasticmapreduce:*","es:*","elasticBeanstalk:*","elasticache:*","eks:*","elasticfilesystem:*","ecs:*","ecr:*","dynamodb:*","dax:*","codeCommit:*","codeBuild:*","cloudWatch:*","cloudTrail:*","cloudSearch:*","cloudFormation:*","backup:*","athena:*","appMesh:*","apiGateway:*","amplify:*","acm:*","ssm:*","rds:*","ec2:*","autoscaling:*","elasticloadbalancing:*","config:*","s3:*","sns:*","logs:*","kms:*","events:*"],"Condition":{"StringLike":{"aws:RequestedRegion":["*"]}},"Effect":"Allow","Resource":"*"},{"Action":["aws-marketplace-management:*","aws-portal:*","ce:*","pricing:*","sts:*","support:*"],"Effect":"Allow","Resource":"*","Sid":"DefaultPermissions"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":["ec2:describeac*","ec2:describeav*","ec2:describesecuritygroups*","ec2:describevpcs*","elasticache:ad*","elasticache:co*","elasticache:createcachec*","elasticache:createcachep*","elasticache:createcachesu*","elasticache:creater*","elasticache:creates*","elasticache:dec*","elasticache:deletecachec*","elasticache:deletecachep*","elasticache:deletecachesu*","elasticache:deleter*","elasticache:deletes*","elasticache:des*","elasticache:i*","elasticache:l*","elasticache:m*","elasticache:p*","elasticache:reb*","elasticache:rem*","elasticache:res*","elasticache:t*","sns:listsubscriptions","sns:listto*"],"Effect":"Allow","Resource":"*","Sid":"AWSelastiCacheOwner"}],"Version":"2012-10-17"}`,

		`{"Statement":[{"Action":["acm:*"],"Effect":"Allow","Resource":"*","Sid":"AWSacmAdmin"}],"Version":"2012-10-17"}`,
		`{"Statement":[{"Action":["es:d*","es:e*","es:g*","es:l*","es:r*","es:s*","es:u*","events:d*","events:e*","events:l*","events:pute*","events:putr*","events:putt*","events:removet*","events:t*","execute-api:*","firehose:*","fsx:*","glacier:abortm*","glacier:ad*","glacier:completem*","glacier:cr*","glacier:d*","glacier:g*","glacier:initiatej*","glacier:initiatem*","glacier:l*","glacier:p*","glacier:r*","glacier:s*","glacier:u*","guardduty:ar*","guardduty:created*","guardduty:createf*","guardduty:createi*","guardduty:creates*","guardduty:createt*","guardduty:dec*","guardduty:deleted*","guardduty:deletef*","guardduty:deleteip*","guardduty:deletet*","guardduty:g*","guardduty:l*","guardduty:s*","guardduty:t*","guardduty:u*","health:describeeventa*","iam:g*","iam:l*","iam:pa*","iam:si*","iam:t*","iam:un*","inspector:*","kafka:*","kinesis:*","kinesisanalytics:*","kinesisvideo:*","kms:*","lambda:*","logs:*","marketplacecommerceanalytics:*","mq:*","pi:*","pricing:*","ram:**","rds:ad*","rds:ap*","rds:b*","rds:co*","rds:createdbc*","rds:createdbi*","rds:createdbp*","rds:createdbsn*","rds:createe*","rds:createo*","rds:deletedbc*","rds:deletedbi*","rds:deletedbp*","rds:deletedbsn*","rds:deletee*","rds:deleteo*","rds:des*","rds:do*","rds:f*","rds:l*","rds:modifyc*","rds:modifydbc*","rds:modifydbi*","rds:modifydbp*","rds:modifydbsn*","rds:modifye*","rds:modifyo*","rds:pr*","rds:reb*","rds:rem*","rds:res*","rds:s*","redshift:ac*","redshift:authorizes*","redshift:b*","redshift:ca*","redshift:co*","redshift:createcluster","redshift:createclusterp*","redshift:createclustersn*","redshift:createclustersu*","redshift:createclusteru*","redshift:createe*","redshift:createsa*","redshift:createsnapshots*","redshift:createt*","redshift:deletecluster","redshift:deleteclusterp*","redshift:deleteclustersn*","redshift:deleteclustersu*","redshift:deletee*","redshift:deletesa*","redshift:deletesnapshots*","redshift:deletet*","redshift:des*","redshift:disables*","redshift:enables*","redshift:ex*","redshift:f*","redshift:g*","redshift:j*","redshift:l*","redshift:m*","redshift:reb*","redshift:res*","redshift:revokes*","redshift:ro*","redshift:v*","robomaker:*","route53:*","route53domains:*","s3:a*","s3:c*","s3:deletea*","s3:deleteb*","s3:deleteo*","s3:g*","s3:h*","s3:l*","s3:puta*","s3:putbuckete*","s3:putbucketn*","s3:putbucketp*","s3:putbuckett*","s3:putbucketv*","s3:putbucketw*","s3:pute*","s3:puti*","s3:putl*","s3:putm*","s3:putobject","s3:putobjectt*","s3:putobjectversiont*","s3:r*","sdb:s*","secretsmanager:*","securityhub:*","shield:*","sns:c*","sns:d*","sns:g*","sns:l*","sns:o*","sns:p*","sns:s*","sns:t*","sns:u*","sqs:c*","sqs:d*","sqs:g*","sqs:l*","sqs:p*","sqs:rec*","sqs:s*","sqs:t*","sqs:u*","ssm:*","states:*","sts:*","support:*","tagging:*","waf-regional:*","waf:*","wafv2:*"],"Effect":"Allow","Resource":"*","Sid":"AWSAdmin2"}],"Version":"2012-10-17"}`,
	}

	for i, testCase := range cases {

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			pol, err := canonicalPolicy(testCase)
			if err != nil {
				t.Errorf("Convert failed for case '%s': %v", pol, err)
			}
			prettyPrint(pol)
		})
	}

}

func prettyPrint(src interface{}) {
	pretty, err := json.MarshalIndent(src, "", "  ")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("\n %s\n", string(pretty))

}
