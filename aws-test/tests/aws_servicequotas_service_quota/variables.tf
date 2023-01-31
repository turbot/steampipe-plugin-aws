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

resource "null_resource" "service_quota" {
  provisioner "local-exec" {
    command = "aws service-quotas get-service-quota --service-code account --quota-code L-77A32B3F > ${path.cwd}/quota.json"
  }
}

data "local_file" "service_quota" {
  depends_on = [null_resource.service_quota]
  filename   = "${path.cwd}/quota.json"
}

output "quota_name" {
  value = jsondecode(data.local_file.service_quota.content).Quota.QuotaName
}

output "quota_code" {
  value = jsondecode(data.local_file.service_quota.content).Quota.QuotaCode
}

output "service_name" {
  value = jsondecode(data.local_file.service_quota.content).Quota.ServiceName
}

output "service_code" {
  value = jsondecode(data.local_file.service_quota.content).Quota.ServiceCode
}

output "resource_aka" {
  value = jsondecode(data.local_file.service_quota.content).Quota.QuotaArn
}
