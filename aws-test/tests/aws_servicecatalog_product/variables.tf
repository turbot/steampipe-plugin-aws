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
  default     = "us-west-2"
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

resource "aws_servicecatalog_product" "named_test_resource" {
  name  = var.resource_name
  owner = var.resource_name
  type  = "CLOUD_FORMATION_TEMPLATE"

  provisioning_artifact_parameters {
    template_url = "https://s3.amazonaws.com/cf-templates-ozkq9d3hgiq2-us-east-1/temp1.json"
  }

  tags = {
    foo = "bar"
  }
}


output "resource_id" {
  value = aws_servicecatalog_product.named_test_resource.id
}

output "resource_aka" {
  value = "arn:aws:catalog:${data.aws_region.primary.name}:${data.aws_caller_identity.current.account_id}:product/${aws_servicecatalog_product.named_test_resource.id}"
}

output "aws_region" {
  depends_on = [aws_servicecatalog_product.named_test_resource]
  value = data.aws_region.primary.name
}

output "aws_account" {
  value = data.aws_caller_identity.current.account_id
}

output "resource_name" {
  value = var.resource_name
}