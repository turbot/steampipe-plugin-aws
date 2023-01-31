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

resource "aws_efs_file_system" "named_test_resource" {
  creation_token   = var.resource_name
  encrypted        = "true"
  performance_mode = "maxIO"
  tags = {
    name = var.resource_name
  }
}

resource "aws_efs_file_system_policy" "efs_file_system_policy" {
  file_system_id = aws_efs_file_system.named_test_resource.id

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Id": "test_policy",
    "Statement": [
        {
            "Sid": "__default_policy_ID",
            "Effect": "Allow",
            "Principal": {
                "AWS": "*"
            },
            "Resource": "${aws_efs_file_system.named_test_resource.arn}",
            "Action": [
                "elasticfilesystem:ClientMount",
                "elasticfilesystem:ClientWrite"
            ],
            "Condition": {
                "Bool": {
                    "aws:SecureTransport": "true"
                }
            }
        }
    ]
}
EOF
}

data "template_file" "resource_aka" {
  template = "arn:$${partition}:elasticfilesystem:$${region}:$${account_id}:file-system/${aws_efs_file_system.named_test_resource.id}"
  vars = {
    partition  = data.aws_partition.current.partition
    account_id = data.aws_caller_identity.current.account_id
    region     = data.aws_region.primary.name
  }
}

output "resource_name" {
  value = var.resource_name
}

output "resource_aka" {
  depends_on = [aws_efs_file_system.named_test_resource]
  value      = data.template_file.resource_aka.rendered
}

output "resource_id" {
  value = aws_efs_file_system.named_test_resource.id
}
