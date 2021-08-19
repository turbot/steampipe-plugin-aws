
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

locals {
  path = "${path.cwd}/output.json"
}

resource "null_resource" "named_test_resource" {
  provisioner "local-exec" {
    command = "aws ec2 describe-reserved-instances --output json > ${local.path}"
  }
}

data "local_file" "input" {
  depends_on = [null_resource.named_test_resource]
  filename   = local.path
}

output "resource_name" {
  depends_on = [null_resource.named_test_resource]
  value      = jsondecode(data.local_file.input.content).ReservedInstances[0].ReservedInstancesId
}

output "instance_type" {
  depends_on = [null_resource.named_test_resource]
  value      = jsondecode(data.local_file.input.content).ReservedInstances[0].InstanceType
}

output "offering_class" {
  depends_on = [null_resource.named_test_resource]
  value      = jsondecode(data.local_file.input.content).ReservedInstances[0].OfferingClass
}

output "state" {
  depends_on = [null_resource.named_test_resource]
  value      = jsondecode(data.local_file.input.content).ReservedInstances[0].State
}

output "scope" {
  depends_on = [null_resource.named_test_resource]
  value      = jsondecode(data.local_file.input.content).ReservedInstances[0].Scope
}

output "usage_price" {
  depends_on = [null_resource.named_test_resource]
  value      = tostring(jsondecode(data.local_file.input.content).ReservedInstances[0].UsagePrice)
}

output "fixed_price" {
  depends_on = [null_resource.named_test_resource]
  value      = tostring(jsondecode(data.local_file.input.content).ReservedInstances[0].FixedPrice)
}

output "currency_code" {
  depends_on = [null_resource.named_test_resource]
  value      = jsondecode(data.local_file.input.content).ReservedInstances[0].CurrencyCode
}


output "resource_aka" {
  depends_on = [null_resource.named_test_resource]
  value = "arn:${data.aws_partition.current.partition}:ec2:${data.aws_region.primary.name}:${data.aws_caller_identity.current.account_id}:instance/${jsondecode(data.local_file.input.content).ReservedInstances[0].ReservedInstancesId}"
}
