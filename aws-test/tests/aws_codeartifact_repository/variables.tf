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

resource "aws_kms_key" "named_test_resource" {
  description = var.resource_name
}

resource "aws_codeartifact_domain" "named_test_resource" {
  domain         = var.resource_name
  encryption_key = aws_kms_key.named_test_resource.arn
  tags = {
    name = var.resource_name
  }
}

resource "aws_codeartifact_repository" "named_test_resource" {
  repository  = var.resource_name
  description = var.resource_name
  domain      = aws_codeartifact_domain.named_test_resource.domain
  tags = {
    name = var.resource_name
  }
}

output "resource_aka" {
  value = aws_codeartifact_repository.named_test_resource.arn
}

output "domain_owner" {
  value = aws_codeartifact_repository.named_test_resource.domain_owner
}

output "administrator_account" {
  value = aws_codeartifact_repository.named_test_resource.administrator_account
}

output "resource_name" {
  value = var.resource_name
}

output "region_name" {
  value = data.aws_region.primary.name
}
