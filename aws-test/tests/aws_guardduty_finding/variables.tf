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

# Create AWS > GuardDuty > Detector
resource "aws_guardduty_detector" "named_test_resource" {
  enable = true
}

locals {
  listPath = "${path.cwd}/listFinding.json"
  getPath  = "${path.cwd}/getFinding.json"
}

resource "null_resource" "named_test_resource" {
  depends_on = [aws_guardduty_detector.named_test_resource]
  provisioner "local-exec" {
    command = "AWS  guardduty create-sample-findings --region ${var.aws_region} --detector-id ${aws_guardduty_detector.named_test_resource.id} --finding-types Backdoor:EC2/DenialOfService.Tcp"
  }
  provisioner "local-exec" {
    command = "aws guardduty list-findings --region ${var.aws_region} --detector-id ${aws_guardduty_detector.named_test_resource.id} > ${local.listPath}"
  }
}

data "local_file" "listInput" {
  depends_on = [null_resource.named_test_resource]
  filename   = local.listPath
}

resource "null_resource" "get_resource" {
  depends_on = [null_resource.named_test_resource]
  provisioner "local-exec" {
    command = "aws guardduty get-findings --region ${var.aws_region} --detector-id ${aws_guardduty_detector.named_test_resource.id} --finding-id ${jsondecode(data.local_file.listInput.content).FindingIds[0]} > ${local.getPath}"
  }
}

data "local_file" "getInput" {
  depends_on = [null_resource.get_resource]
  filename   = local.getPath
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "resource_name" {
  depends_on = [null_resource.get_resource]
  value      = jsondecode(data.local_file.getInput.content).Findings[0].Title
}

output "resource_aka" {
  depends_on = [null_resource.get_resource]
  value      = jsondecode(data.local_file.getInput.content).Findings[0].Arn
}

output "resource_id" {
  depends_on = [null_resource.get_resource]
  value      = jsondecode(data.local_file.getInput.content).Findings[0].Id
}
