variable "resource_name" {
  type        = string
  default     = ""
  description = "Name of the resource used throughout the test."
}

variable "turbot_profile" {
  type        = string
  default     = ""
  description = "Turbot credentials profile to use for the test run."
}

provider "turbot" {
  profile = var.turbot_profile
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

resource "null_resource" "delay" {
  provisioner "local-exec" {
    command = "sleep 300"
  }
}

resource "null_resource" "named_test_resource" {
  depends_on = [null_resource.delay]
  provisioner "local-exec" {
    command = <<EOT
      aws redshift delete-cluster-snapshot --snapshot-identifier ${var.resource_name} --profile ${var.aws_profile} --region ${var.aws_region};
    EOT
  }
}
