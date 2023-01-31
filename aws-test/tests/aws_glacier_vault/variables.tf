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

resource "aws_glacier_vault" "named_test_resource" {
  name          = var.resource_name
  access_policy = <<EOF
{
    "Version":"2012-10-17",
    "Statement":[
       {
          "Sid": "__default_policy_ID",
          "Principal": "*",
          "Effect": "Allow",
          "Action": [
             "glacier:InitiateJob",
             "glacier:GetJobOutput"
          ],
          "Resource": "arn:${data.aws_partition.current.partition}:glacier:${data.aws_region.primary.name}:${data.aws_caller_identity.current.account_id}:vaults/${var.resource_name}"
       }
    ]
}
EOF
  tags = {
    name = var.resource_name
  }
}

resource "aws_glacier_vault_lock" "vault_policy" {
  complete_lock = false
  policy        = <<EOF
{
    "Version":"2012-10-17",
    "Statement":[
       {
          "Sid": "deny-based-on-archive-age",
          "Principal": "*",
          "Effect": "Deny",
          "Action": [
             "glacier:DeleteArchive"
          ],
          "Resource": "${aws_glacier_vault.named_test_resource.arn}",
          "Condition": {
            "NumericLessThan": {
              "glacier:ArchiveAgeInDays" : "365"
            }
          }
       }
    ]
}
EOF
  vault_name    = aws_glacier_vault.named_test_resource.name
}

output "resource_name" {
  value = var.resource_name
}

output "resource_aka" {
  value = aws_glacier_vault.named_test_resource.arn
}

output "resource_id" {
  value = aws_glacier_vault.named_test_resource.name
}

