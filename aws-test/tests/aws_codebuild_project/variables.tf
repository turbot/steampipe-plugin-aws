
variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125"
  description = "Name of the resource used throughout the test."
}

variable "aws_profile" {
  type        = string
  default     = "default"
  description = "AWS credentials profile used for the test. Default is to use the default profile."
}

variable "aws_region" {
  type        = string
  default     = "us-east-2"
  description = "AWS region used for the test. Does not work with default region in config, so must be defined here."
}

variable "aws_region_alternate" {
  type        = string
  default     = "us-east-2"
  description = "Alternate AWS region used for tests that require two regions."
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

resource "aws_vpc" "my_vpc" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_subnet" "my_subnet1" {
  vpc_id     = aws_vpc.my_vpc.id
  cidr_block = "10.1.1.0/24"
}

resource "aws_subnet" "my_subnet2" {
  vpc_id     = aws_vpc.my_vpc.id
  cidr_block = "10.1.2.0/24"
}

resource "aws_s3_bucket" "test_bucket" {
  bucket        = var.resource_name
  force_destroy = true
}

resource "aws_iam_role" "test_role" {
  name               = var.resource_name
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "codebuild.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "test_role_policy" {
  role   = aws_iam_role.test_role.name
  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Resource": [
        "*"
      ],
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "ec2:CreateNetworkInterface",
        "ec2:DescribeDhcpOptions",
        "ec2:DescribeNetworkInterfaces",
        "ec2:DeleteNetworkInterface",
        "ec2:DescribeSubnets",
        "ec2:DescribeSecurityGroups",
        "ec2:DescribeVpcs"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "ec2:CreateNetworkInterfacePermission"
      ],
      "Resource": [
        "arn:*:ec2:::network-interface/*"
      ],
      "Condition": {
        "StringEquals": {
          "ec2:Subnet": [
            "${aws_subnet.my_subnet1.arn}",
            "${aws_subnet.my_subnet2.arn}"
          ],
          "ec2:AuthorizedService": "codebuild.amazonaws.com"
        }
      }
    },
    {
      "Effect": "Allow",
      "Action": [
        "s3:*"
      ],
      "Resource": [
        "${aws_s3_bucket.test_bucket.arn}",
        "${aws_s3_bucket.test_bucket.arn}/*"
      ]
    }
  ]
}
POLICY
}

resource "aws_codebuild_project" "named_test_resource" {
  name         = var.resource_name
  description  = "Resource for testing"
  service_role = aws_iam_role.test_role.arn
  artifacts {
    type = "NO_ARTIFACTS"
  }
  cache {
    type     = "S3"
    location = aws_s3_bucket.test_bucket.bucket
  }
  environment {
    compute_type                = "BUILD_GENERAL1_SMALL"
    image                       = "aws/codebuild/standard:1.0"
    type                        = "LINUX_CONTAINER"
    image_pull_credentials_type = "CODEBUILD"
  }
  source {
    type     = "GITHUB"
    location = "https://github.com/mitchellh/packer.git"
  }
  tags = {
    foo = "bar"
  }
}

output "resource_aka" {
  value = aws_codebuild_project.named_test_resource.arn
}

output "resource_id" {
  value = aws_codebuild_project.named_test_resource.id
}

output "resource_name" {
  value = var.resource_name
}
