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

resource "aws_vpc" "my_vpc" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "my_subnet" {
  vpc_id     = aws_vpc.my_vpc.id
  cidr_block = "10.0.0.0/24"
}

resource "aws_elasticache_subnet_group" "my_subnet_group" {
  name       = var.resource_name
  subnet_ids = [aws_subnet.my_subnet.id]
}

resource "aws_elasticache_replication_group" "named_test_resource" {
  replication_group_id          = var.resource_name
  automatic_failover_enabled    = true
  replication_group_description = "test description"
  node_type                     = "cache.t2.micro"
  num_cache_clusters            = 2
  parameter_group_name          = "default.redis5.0"
  engine_version                = "5.0.6"
  port                          = 6379
  subnet_group_name             = aws_elasticache_subnet_group.my_subnet_group.id
}

output "resource_aka" {
  value = aws_elasticache_replication_group.named_test_resource.arn
}

output "resource_id" {
  value = aws_elasticache_replication_group.named_test_resource.id
}