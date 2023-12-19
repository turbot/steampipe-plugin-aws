variable "resource_name" {
  type        = string
  default     = "turbot-test-20210612"
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

resource "aws_ec2_managed_prefix_list" "named_test_resource" {
  name           = var.resource_name
  address_family = "IPv4"
  max_entries    = 5
}

resource "aws_vpc" "named_test_resource" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_ec2_managed_prefix_list_entry" "entry_1" {
  cidr           = aws_vpc.named_test_resource.cidr_block
  description    = "Primary"
  prefix_list_id = aws_ec2_managed_prefix_list.named_test_resource.id
}

output "address_family" {
  value = aws_ec2_managed_prefix_list.named_test_resource.address_family
}

output "cidr_block" {
  value = aws_vpc.named_test_resource.cidr_block
}

output "prefix_list_id" {
  value = aws_ec2_managed_prefix_list.named_test_resource.id
}

output "owner_id" {
  value = aws_ec2_managed_prefix_list.named_test_resource.owner_id
}

output "resource_name" {
  value = var.resource_name
}
