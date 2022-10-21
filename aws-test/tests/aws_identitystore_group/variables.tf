variable "resource_name" {
  type        = string
  default     = "turbot-test-20200155-create-update"
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
  default     = "us-east-1"
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

data "aws_ssoadmin_instances" "main" {}

resource "aws_identitystore_group" "main" {
  display_name      = var.resource_name
  description       = "Example description"
  identity_store_id = tolist(data.aws_ssoadmin_instances.main.identity_store_ids)[0]
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "aws_region" {
  value = data.aws_region.primary.name
}

output "aws_partition" {
  value = data.aws_partition.current.partition
}

output "resource_id" {
  value = resource.aws_identitystore_group.main.group_id
}

output "resource_name" {
  value = resource.aws_identitystore_group.main.display_name
}

output "identity_store_id" {
  value = resource.aws_identitystore_group.main.identity_store_id
}
