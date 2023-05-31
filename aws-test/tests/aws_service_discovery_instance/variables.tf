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

resource "aws_vpc" "named_test_resource" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true
}

resource "aws_service_discovery_private_dns_namespace" "named_test_resource" {
  name        = var.resource_name
  description = "example"
  vpc         = aws_vpc.named_test_resource.id
}

resource "aws_service_discovery_service" "named_test_resource" {
  name = var.resource_name

  dns_config {
    namespace_id = aws_service_discovery_private_dns_namespace.named_test_resource.id

    dns_records {
      ttl  = 10
      type = "A"
    }

    routing_policy = "MULTIVALUE"
  }

  health_check_custom_config {
    failure_threshold = 1
  }

  tags = {
    name = var.resource_name
  }
}

resource "aws_service_discovery_instance" "named_test_resource" {
  instance_id = var.resource_name
  service_id  = aws_service_discovery_service.named_test_resource.id

  attributes = {
    AWS_INSTANCE_IPV4 = "172.18.0.1"
    custom_attribute  = "custom"
  }
}

output "resource_id" {
  value = aws_service_discovery_instance.named_test_resource.id
}

output "service_id" {
  value = aws_service_discovery_service.named_test_resource.id
}

output "aws_region" {
  value = data.aws_region.primary.name
}

output "aws_account" {
  value = data.aws_caller_identity.current.account_id
}

output "resource_name" {
  value = var.resource_name
}