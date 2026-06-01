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
  description = "Alternate AWS region used for tests that require two regions (e.g. ECR replication)."
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

# `aws_region.name` is deprecated in newer AWS provider versions; use `region`
# in resource ARNs / outputs and rely on these locals so the file works on both
# old and new provider versions.
locals {
  primary_region   = data.aws_region.primary.region
  alternate_region = data.aws_region.alternate.region
}

# Configure cross-region replication on the private registry. There is exactly
# one registry per (account, region), so this resource is referenced by the
# table directly.
resource "aws_ecr_replication_configuration" "test" {
  replication_configuration {
    rule {
      destination {
        region      = var.aws_region_alternate
        registry_id = data.aws_caller_identity.current.account_id
      }
    }
  }
}

# Attach a registry permissions policy so the `policy` and `policy_std`
# columns can be exercised.
resource "aws_ecr_registry_policy" "test" {
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "turbot-test-allow-self-replicate"
        Effect    = "Allow"
        Principal = {
          AWS = "arn:${data.aws_partition.current.partition}:iam::${data.aws_caller_identity.current.account_id}:root"
        }
        Action = [
          "ecr:ReplicateImage",
          "ecr:CreateRepository"
        ]
        Resource = "arn:${data.aws_partition.current.partition}:ecr:${local.primary_region}:${data.aws_caller_identity.current.account_id}:repository/*"
      }
    ]
  })
}

output "registry_id" {
  value = data.aws_caller_identity.current.account_id
}

output "aws_account" {
  value = data.aws_caller_identity.current.account_id
}

output "aws_region" {
  value = local.primary_region
}

output "aws_region_alternate" {
  value = local.alternate_region
}

output "aws_partition" {
  value = data.aws_partition.current.partition
}

output "resource_aka" {
  value = "arn:${data.aws_partition.current.partition}:ecr:${local.primary_region}:${data.aws_caller_identity.current.account_id}:registry/${data.aws_caller_identity.current.account_id}"
}
