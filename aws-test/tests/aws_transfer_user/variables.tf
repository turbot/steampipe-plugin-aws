

variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update"
  description = "Name of the resource used throughout the test."
}

variable "identity_provider_type" {
  type        = string
  default     = "SERVICE_MANAGED"
  description = "Name of the identity provider type."
}

variable "domain" {
  type        = string
  default     = "S3"
  description = "Domain of the storage system used for file transfers."
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

data "aws_canonical_user_id" "current_user" {}
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

resource "aws_transfer_server" "resource" {
  domain = var.domain
  identity_provider_type = var.identity_provider_type
  tags = {
    Name = var.resource_name
  }
}

resource "aws_iam_role" "resource" {
  name               = "tf-test-transfer-user-iam-role"
  assume_role_policy = jsonencode({
    statement {
      effect = "Allow"

      principals {
        type        = "Service"
        identifiers = ["transfer.amazonaws.com"]
      }

      actions = ["sts:AssumeRole"]
    }
  })
}

resource "aws_transfer_user" "named_test_resource" {
  server_id = aws_transfer_server.resource.id
  user_name = var.resource_name
  role      = aws_iam_role.resource.arn
}

output "canonical_user_id" {
  value = data.aws_canonical_user_id.current_user.id
}

output "resource_server_id" {
  value = aws_transfer_user.named_test_resource.server_id
}

output "resource_aka" {
  value = aws_transfer_user.named_test_resource.arn
}

output "resource_user_name" {
  value = aws_transfer_user.named_test_resource.user_name
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "aws_partition" {
  value = data.aws_partition.current.partition
}
