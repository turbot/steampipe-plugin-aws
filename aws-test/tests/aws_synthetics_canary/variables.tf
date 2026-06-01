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

resource "local_file" "javascript_file" {
  filename          = "${path.cwd}/../../test.js"
  sensitive_content = "module.exports = { handler: () => { console.log('This is an integration test for creation of canaries') } }"
}

data "archive_file" "zip" {
  depends_on  = [local_file.javascript_file]
  type        = "zip"
  source_file = "${path.cwd}/../../test.js"
  output_path = "${path.cwd}/../../test.zip"
}

resource "aws_s3_bucket" "artifacts" {
  bucket = var.resource_name
}

resource "aws_iam_role" "aws_synthetics_canary" {
  name = var.resource_name
  assume_role_policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Action" : "sts:AssumeRole",
        "Principal" : {
          "Service" : "lambda.amazonaws.com"
        },
        "Effect" : "Allow",
        "Sid" : ""
      }
    ]
  })
}

resource "aws_synthetics_canary" "named_test_resource" {
  name                 = var.resource_name
  artifact_s3_location = "s3://${var.resource_name}"
  zip_file             = data.archive_file.zip.output_path
  execution_role_arn   = aws_iam_role.aws_synthetics_canary.arn
  handler              = "exports.handler"
  runtime_version      = "syn-nodejs-puppeteer-15.0"

  schedule {
    expression = "rate(1 hour)"
  }
}

output "resource_aka" {
  value = aws_synthetics_canary.named_test_resource.id
}

output "resource_id" {
  value = aws_synthetics_canary.named_test_resource.name
}

output "artifact_s3_location" {
  value = var.resource_name
}

output "code_handler" {
  value = aws_synthetics_canary.named_test_resource.handler
}

output "schedule_expression" {
  value = aws_synthetics_canary.named_test_resource.schedule[0].expression
}

output "iam_role_arn" {
  value = aws_iam_role.aws_synthetics_canary.arn
}

output "resource_name" {
  value = var.resource_name
}

output "runtime_version" {
  value = aws_synthetics_canary.named_test_resource.runtime_version
}

output "aws_partition" {
  value = data.aws_partition.current.partition
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "aws_region" {
  value = data.aws_region.primary.name
}

output "canary_arn" {
  value = aws_synthetics_canary.named_test_resource.arn
}
