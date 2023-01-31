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

resource "aws_backup_vault" "named_test_resource" {
  name = var.resource_name
}

resource "aws_backup_plan" "named_test_resource" {
  name = var.resource_name

  rule {
    rule_name         = var.resource_name
    target_vault_name = aws_backup_vault.named_test_resource.name
    schedule          = "cron(0 12 * * ? *)"
  }
}

resource "aws_iam_role" "example" {
  name               = var.resource_name
  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": ["sts:AssumeRole"],
      "Effect": "allow",
      "Principal": {
        "Service": ["backup.amazonaws.com"]
      }
    }
  ]
}
POLICY
}

resource "aws_iam_role_policy_attachment" "example" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSBackupServiceRolePolicyForBackup"
  role       = aws_iam_role.example.name
}

resource "aws_ebs_volume" "named_test_resource" {
  availability_zone = "${var.aws_region}a"
  size              = 1
  tags = {
    Name = var.resource_name
  }
}

resource "aws_backup_selection" "named_test_resource" {
  iam_role_arn = aws_iam_role.example.arn
  name         = var.resource_name
  plan_id      = aws_backup_plan.named_test_resource.id

  resources = [
    aws_ebs_volume.named_test_resource.arn,
  ]
}

output "plan_id" {
  value = aws_backup_plan.named_test_resource.id
}

output "selection_id" {
  value = aws_backup_selection.named_test_resource.id
}

output "resource_aka" {
  value = "arn:${data.aws_partition.current.partition}:backup:${var.aws_region}:${data.aws_caller_identity.current.account_id}:backup-plan:${aws_backup_plan.named_test_resource.id}/selection/${aws_backup_selection.named_test_resource.id}"
}

output "volume_arn" {
  value = aws_ebs_volume.named_test_resource.arn
}

output "iam_role_arn" {
  value = aws_iam_role.example.arn
}

output "resource_name" {
  value = var.resource_name
}
