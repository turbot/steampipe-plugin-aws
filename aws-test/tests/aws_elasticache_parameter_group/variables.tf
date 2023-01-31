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

resource "aws_elasticache_parameter_group" "named_test_resource" {
  name        = var.resource_name
  family      = "redis2.8"
  description = "A test parameter group"
}

data "template_file" "resource_aka" {
  template = "arn:$${partition}:elasticache:$${region}:$${account_id}:parametergroup:${aws_elasticache_parameter_group.named_test_resource.name}"
  vars = {
    resource_name    = var.resource_name
    partition        = data.aws_partition.current.partition
    account_id       = data.aws_caller_identity.current.account_id
    region           = data.aws_region.primary.name
    alternate_region = data.aws_region.alternate.name
  }
}

output "resource_name" {
  value = var.resource_name
}

output "resource_aka" {
  depends_on = [aws_elasticache_parameter_group.named_test_resource]
  value      = data.template_file.resource_aka.rendered
}
