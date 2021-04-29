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

resource "aws_s3_bucket" "named_test_resource" {
  bucket = var.resource_name
}

resource "aws_vpc" "named_test_resource" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_s3_access_point" "named_test_resource" {
  bucket = aws_s3_bucket.named_test_resource.id
  name   = var.resource_name

  vpc_configuration {
    vpc_id = aws_vpc.named_test_resource.id
  }
}

output "vpc_id" {
  value = aws_vpc.named_test_resource.id
}

output "resource_id" {
  value = aws_s3_access_point.named_test_resource.id
}

output "arn" {
  value = aws_s3_access_point.named_test_resource.arn
}

output "network_origin" {
  value = aws_s3_access_point.named_test_resource.network_origin
}

output "has_public_access_policy" {
  value = aws_s3_access_point.named_test_resource.has_public_access_policy
}

output "resource_name" {
  value = var.resource_name
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "aws_partition" {
  value = data.aws_partition.current.partition
}

output "region_id" {
  value = split(":", aws_s3_access_point.named_test_resource.arn)[3]
}
