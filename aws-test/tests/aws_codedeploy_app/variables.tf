variable "resource_name" {
  type        = string
  default     = "turbot-test-20200126"
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

resource "aws_codedeploy_app" "named_test_resource" {
  compute_platform = "Lambda"
  name             = var.resource_name
  tags = {
    name = var.resource_name
  }
}

output "resource_aka" {
  value = aws_codedeploy_app.named_test_resource.arn
}

output "resource_id" {
  value = aws_codedeploy_app.named_test_resource.application_id
}

output "compute_platform" {
  value = aws_codedeploy_app.named_test_resource.compute_platform
}

output "linked_to_github" {
  value = aws_codedeploy_app.named_test_resource.linked_to_github
}

output "github_account_name" {
  value = aws_codedeploy_app.named_test_resource.github_account_name
}

output "resource_name" {
  value = var.resource_name
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}
