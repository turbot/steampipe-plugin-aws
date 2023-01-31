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
  default     = "us-east-2"
  description = "AWS region used for the test. Does not work with default region in config, so must be defined here."
}

variable "aws_region_alternate" {
  type        = string
  default     = "us-east-2"
  description = "Alternate AWS region used for tests that require two regions (e.g. DynamoDB global tables)."
}

variable "description" {
  type        = string
  default     = "integration testing"
  description = "Description used for the test."
}

variable "connection_type" {
  type        = string
  default     = "JDBC"
  description = "connection type used for the test."
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

resource "aws_glue_connection" "named_test_resource" {
  name            = var.resource_name
  description     = var.description
  connection_type = var.connection_type
  connection_properties = {
    JDBC_CONNECTION_URL = "jdbc:mysql://example.com/exampledatabase"
    PASSWORD            = "examplepassword"
    USERNAME            = "exampleusername"
  }
}


output "id" {
  value = aws_glue_connection.named_test_resource.id
}

output "resource_aka" {
  value = aws_glue_connection.named_test_resource.arn
}

output "description" {
  value = var.description
}

output "connection_type" {
  value = var.connection_type
}

output "aws_region" {
  value = data.aws_region.primary.name
}

output "resource_name" {
  value = var.resource_name
}

output "aws_account" {
  value = data.aws_caller_identity.current.account_id
}
