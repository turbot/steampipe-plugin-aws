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

resource "aws_kinesis_stream" "named_test_resource" {
  name                      = var.resource_name
  shard_count               = 1
  retention_period          = 24
  enforce_consumer_deletion = true
}

resource "null_resource" "named_test_resource" {
  depends_on = [aws_kinesis_stream.named_test_resource]
  provisioner "local-exec" {
    command = "aws kinesis register-stream-consumer --stream-arn ${aws_kinesis_stream.named_test_resource.arn} --consumer-name ${var.resource_name} --profile ${var.aws_profile} --region  ${data.aws_region.primary.name}"
  }
}

locals {
  path = "${path.cwd}/kinesis_stream.json"
}

resource "null_resource" "output_test_resource" {
  depends_on = [null_resource.named_test_resource]
  provisioner "local-exec" {
    command = "aws kinesis list-stream-consumers --stream-arn ${aws_kinesis_stream.named_test_resource.arn} --output json --profile ${var.aws_profile} --region ${data.aws_region.primary.name} > ${local.path}"
  }
}

data "local_file" "input" {
  depends_on = [null_resource.output_test_resource]
  filename   = local.path
}

output "resource_aka" {
  depends_on = [null_resource.output_test_resource]
  value      = jsondecode(data.local_file.input.content).Consumers[0].ConsumerARN
}

output "stream_aka" {
  value = aws_kinesis_stream.named_test_resource.arn
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "aws_region" {
  value = data.aws_region.primary.name
}

output "resource_name" {
  value = var.resource_name
}
