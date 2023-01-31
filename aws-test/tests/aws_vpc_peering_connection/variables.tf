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

resource "aws_vpc" "first_vpc" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_vpc" "second_vpc" {
  depends_on = [aws_vpc.first_vpc]
  cidr_block = "10.2.0.0/16"
}

resource "aws_vpc_peering_connection" "named_test_resource" {
  depends_on  = [aws_vpc.second_vpc]
  peer_vpc_id = aws_vpc.second_vpc.id
  vpc_id      = aws_vpc.first_vpc.id
  auto_accept = true

  tags = {
    name = var.resource_name
  }
}

output "resource_name" {
  value = var.resource_name
}

output "id" {
  value = aws_vpc_peering_connection.named_test_resource.id
}

output "peer_vpc_id" {
  value = aws_vpc_peering_connection.named_test_resource.peer_vpc_id
}

output "vpc_id" {
  value = aws_vpc_peering_connection.named_test_resource.vpc_id
}
