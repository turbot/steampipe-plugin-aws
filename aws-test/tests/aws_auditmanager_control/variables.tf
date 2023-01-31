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

resource "null_resource" "named_test_resource" {
  provisioner "local-exec" {
    command = <<EOT
      aws auditmanager create-control --name ${var.resource_name} --control-mapping-sources "sourceName"=${var.resource_name},"sourceType"="AWS_Cloudtrail","sourceSetUpOption"="System_Controls_Mapping",sourceKeyword="{keywordInputType="SELECT_FROM_LIST",keywordValue="a4b_ApproveSkill"}" --tags name=${var.resource_name} --profile ${var.aws_profile} > ${path.cwd}/control.json;
    EOT
  }
}

data "local_file" "control" {
  depends_on = [null_resource.named_test_resource]
  filename   = "${path.cwd}/control.json"
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "aws_region" {
  value = data.aws_region.primary.name
}

output "aws_partition" {
  value = data.aws_partition.current.partition
}

output "control_id" {
  value = jsondecode(data.local_file.control.content).control.id
}

output "resource_aka" {
  value = jsondecode(data.local_file.control.content).control.arn
}

output "resource_name" {
  value = var.resource_name
}
