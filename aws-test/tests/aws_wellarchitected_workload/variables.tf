variable "resource_name" {
  type    = string
  default = ""
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

locals {
  path = "${path.cwd}/kinesis_stream.json"
}

resource "null_resource" "named_test_resource" {
  provisioner "local-exec" {
    command = "aws wellarchitected create-workload --workload-name ${var.resource_name} --description 'testing-plugin'  --environment 'PREPRODUCTION' --review-owner 'abc' --lenses 'wellarchitected' --aws-regions 'us-east-1' > ${local.path}"
  }
}

data "local_file" "input" {
  depends_on = [null_resource.named_test_resource]
  filename   = local.path
}

output "resource_id" {
  depends_on = [null_resource.named_test_resource]
  value = jsondecode(data.local_file.input.content).WorkloadId
}

output "resource_name" {
  value = var.resource_name
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "aws_region" {
  value = data.aws_region.primary.name
}

data "template_file" "resource_aka" {
  template = "arn:$${partition}:wellarchitected:$${region}:$${account_id}:workload/$${resource_id}"
  vars = {
    partition = data.aws_partition.current.partition
    account_id = data.aws_caller_identity.current.account_id
    region = data.aws_region.primary.name
    resource_id = jsondecode(data.local_file.input.content).WorkloadId
    alternate_region = data.aws_region.alternate.name
  }
}

output "resource_aka" {
  depends_on = [ null_resource.named_test_resource ]
  value = data.template_file.resource_aka.rendered
}

output "aws_partition" {
  value = data.aws_partition.current.partition
}
