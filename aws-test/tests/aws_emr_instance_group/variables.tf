variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update"
  description = "Name of the resource used throughout the test."
}

variable "aws_profile" {
  type        = string
  default     = "default"
  description = "AWS credentials profile used for the test. Default is to use the default profile."
}

variable "aws_region" {
  type        = string
  default     = "us-east-1"
  description = "AWS region used for the test. Does not work with default region in config, so must be defined here."
}

variable "aws_region_alternate" {
  type        = string
  default     = "us-east-2"
  description = "Alternate AWS region used for tests that require two regions (e.g. DynamoDB global tables)."
}

provider "aws" {
  profile = var.aws_profile
  region  = var.aws_region
}

provider "aws" {
  alias   = "alternate"
  profile = var.aws_profile
  region  = var.aws_region_alternate
}

data "aws_partition" "current" {}
data "aws_caller_identity" "current" {}
data "aws_region" "primary" {}
data "aws_region" "alternate" {
  provider = aws.alternate
}

data "null_data_source" "resource" {
  inputs = {
    scope = "arn:${data.aws_partition.current.partition}:::${data.aws_caller_identity.current.account_id}"
  }
}

resource "aws_iam_role" "iam_emr_service_role" {
  name               = "${var.resource_name}_1"
  assume_role_policy = <<EOF
{
  "Version": "2008-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "elasticmapreduce.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "iam_emr_service_policy" {
  name   = "${var.resource_name}_1"
  role   = aws_iam_role.iam_emr_service_role.id
  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [{
        "Effect": "Allow",
        "Resource": "*",
        "Action": [
            "ec2:AuthorizeSecurityGroupEgress",
            "ec2:AuthorizeSecurityGroupIngress",
            "ec2:CancelSpotInstanceRequests",
            "ec2:CreateNetworkInterface",
            "ec2:CreateSecurityGroup",
            "ec2:CreateTags",
            "ec2:DeleteNetworkInterface",
            "ec2:DeleteSecurityGroup",
            "ec2:DeleteTags",
            "ec2:DescribeAvailabilityZones",
            "ec2:DescribeAccountAttributes",
            "ec2:DescribeDhcpOptions",
            "ec2:DescribeInstanceStatus",
            "ec2:DescribeInstances",
            "ec2:DescribeKeyPairs",
            "ec2:DescribeNetworkAcls",
            "ec2:DescribeNetworkInterfaces",
            "ec2:DescribePrefixLists",
            "ec2:DescribeRouteTables",
            "ec2:DescribeSecurityGroups",
            "ec2:DescribeSpotInstanceRequests",
            "ec2:DescribeSpotPriceHistory",
            "ec2:DescribeSubnets",
            "ec2:DescribeVpcAttribute",
            "ec2:DescribeVpcEndpoints",
            "ec2:DescribeVpcEndpointServices",
            "ec2:DescribeVpcs",
            "ec2:DetachNetworkInterface",
            "ec2:ModifyImageAttribute",
            "ec2:ModifyInstanceAttribute",
            "ec2:RequestSpotInstances",
            "ec2:RevokeSecurityGroupEgress",
            "ec2:RunInstances",
            "ec2:TerminateInstances",
            "ec2:DeleteVolume",
            "ec2:DescribeVolumeStatus",
            "ec2:DescribeVolumes",
            "ec2:DetachVolume",
            "iam:GetRole",
            "iam:GetRolePolicy",
            "iam:ListInstanceProfiles",
            "iam:ListRolePolicies",
            "iam:PassRole",
            "s3:CreateBucket",
            "s3:Get*",
            "s3:List*",
            "sdb:BatchPutAttributes",
            "sdb:Select",
            "sqs:CreateQueue",
            "sqs:Delete*",
            "sqs:GetQueue*",
            "sqs:PurgeQueue",
            "sqs:ReceiveMessage"
        ]
    }]
}
EOF
}

resource "aws_iam_role" "iam_emr_profile_role" {
  name               = "${var.resource_name}_2"
  assume_role_policy = <<EOF
{
  "Version": "2008-10-17",
  "Statement": [
    {
      "Sid": "test",
      "Effect": "Allow",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_instance_profile" "emr_profile" {
  name = var.resource_name
  role = aws_iam_role.iam_emr_profile_role.name
}

resource "aws_iam_role_policy" "iam_emr_profile_policy" {
  name   = "${var.resource_name}_2"
  role   = aws_iam_role.iam_emr_profile_role.id
  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [{
        "Effect": "Allow",
        "Resource": "*",
        "Action": [
            "cloudwatch:*",
            "ec2:Describe*",
            "elasticmapreduce:Describe*",
            "elasticmapreduce:ListBootstrapActions",
            "elasticmapreduce:ListClusters",
            "elasticmapreduce:ListInstanceGroups",
            "elasticmapreduce:ListInstances",
            "elasticmapreduce:ListSteps",
            "s3:*",
            "sdb:*",
            "sns:*",
            "sqs:*"
        ]
    }]
}
EOF
}

resource "aws_vpc" "my_vpc" {
  cidr_block = "168.31.0.0/16"
}

resource "aws_security_group" "emr_master" {
  vpc_id                 = aws_vpc.my_vpc.id
  revoke_rules_on_delete = true
  ingress {
    from_port = 443
    to_port   = 443
    protocol  = "tcp"
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
  depends_on = [aws_subnet.my_subnet]
  lifecycle {
    ignore_changes = ["ingress", "egress"]
  }
}

resource "aws_subnet" "my_subnet" {
  vpc_id     = aws_vpc.my_vpc.id
  cidr_block = "168.31.0.0/20"
}

resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.my_vpc.id
}

resource "aws_route_table" "test_route_table" {
  vpc_id = aws_vpc.my_vpc.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.gw.id
  }
}

resource "aws_main_route_table_association" "test_association" {
  vpc_id         = aws_vpc.my_vpc.id
  route_table_id = aws_route_table.test_route_table.id
}

resource "aws_key_pair" "deployer" {
  key_name   = var.resource_name
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC6B60dQX19dOO0wmNnPf27NJBXQhuVAsFcshY7YqTmNA1dfl0lXvJFSF0nVTvaxtE2854gmvygm/yafopG33zJeH+k2ZrblivkEuH0NNajKMp2mBumkk8RiTVvU1ET4m6Z1jZaK0dTSeV2k6gUk7Wmvo4+6RIu+wdAtwGGGcSq4HJ4M0G1CGHruBywMdOs5CklqRn7BRMh2aexowDachlabUKRFwTf7OCjDcoCWHo+o4kHV7NDEdhFzrE2fxXttt054kd5CMZOEOB9p01pAqSPrxFv3S/BpiXuhmjwsUsIdfHIgFrFhNY/M4+KhTa3pI69NqsigReotTEzk2QwHyHSjsRwfzeJNRgI8iqGIdYr3+t6zCKswyzgV0hqalZhXtpVUEBbumJAKWy+VzTnOXUQ5zEVvabTxi3dsW9zMGofjwhiC75x97RE7U/ZEanRM+0Avr8e3c3ob4NLf9JOwbvB/jYvS6j/nuwJ5H6igT3Oj0oDiknh2WtOgH/BKKaCq5Nyl15Df65i0PNgKRcLy+x0fubYUtsCOMPpQTfztSYGjbvAhAOrH2LJF0i2dTD3PcVw6R6uDpJSw9QkzdD+2ofxnqjGPUeKR3mraGOXt2hZe9F9bvvdZScDfdAoyBUYTO2/sxsQFobRtdiEk2up+mNIx9QXMXO2dOgtNUYZKc1DdQ== clara@turbot.com"
}

resource "aws_emr_cluster" "named_test_resource" {
  name          = var.resource_name
  release_label = "emr-4.9.5"
  applications  = ["Spark"]
  ec2_attributes {
    key_name                          = aws_key_pair.deployer.key_name
    subnet_id                         = aws_subnet.my_subnet.id
    emr_managed_master_security_group = aws_security_group.emr_master.id
    emr_managed_slave_security_group  = aws_security_group.emr_master.id
    instance_profile                  = aws_iam_instance_profile.emr_profile.arn
  }
  master_instance_group {
    instance_type = "m4.large"
  }
  service_role = aws_iam_role.iam_emr_service_role.arn
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = aws_emr_cluster.named_test_resource.id
}

output "region_name" {
  value = data.aws_region.primary.name
}
