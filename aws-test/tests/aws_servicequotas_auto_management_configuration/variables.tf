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

resource "null_resource" "auto_management_configuration" {
  provisioner "local-exec" {
    command = "aws service-quotas get-auto-management-configuration > ${path.cwd}/auto_management_configuration.json"
  }
}

data "local_file" "auto_management_configuration" {
  depends_on = [null_resource.auto_management_configuration]
  filename   = "${path.cwd}/auto_management_configuration.json"
}

output "opt_in_status" {
  value = jsondecode(data.local_file.auto_management_configuration.content).AutoManagementConfiguration.OptInStatus
}

output "opt_in_type" {
  value = jsondecode(data.local_file.auto_management_configuration.content).AutoManagementConfiguration.OptInType
}

output "opt_in_level" {
  value = jsondecode(data.local_file.auto_management_configuration.content).AutoManagementConfiguration.OptInLevel
}

output "notification_arn" {
  value = jsondecode(data.local_file.auto_management_configuration.content).AutoManagementConfiguration.NotificationArn
}

output "exclusion_list" {
  value = jsondecode(data.local_file.auto_management_configuration.content).AutoManagementConfiguration.ExclusionList
}

output "resource_aka" {
  value = "arn:${data.aws_partition.current.partition}:servicequotas:${data.aws_region.primary.name}:${data.aws_caller_identity.current.account_id}:auto-management-configuration"
}
