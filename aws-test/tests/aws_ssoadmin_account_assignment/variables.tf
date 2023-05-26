variable "aws_profile" {
  type        = string
  default     = null
  description = "AWS credentials profile used for the test. Default is to use AWS_PROFILE environment variable, then default profile."
}

provider "aws" {
  profile = var.aws_profile
}

data "aws_partition" "current" {}
data "aws_caller_identity" "current" {}
data "aws_region" "primary" {}

data "aws_ssoadmin_instances" "main" {}

resource "aws_ssoadmin_permission_set" "main" {
  name             = "steampipe-test"
  instance_arn     = one(data.aws_ssoadmin_instances.main.arns)
}

resource "aws_identitystore_group" "main" {
  display_name      = "steampipe-test"
  description       = "Example description"
  identity_store_id = one(data.aws_ssoadmin_instances.main.identity_store_ids)
}

resource "aws_ssoadmin_account_assignment" "main" {
  instance_arn       = one(data.aws_ssoadmin_instances.main.arns)
  permission_set_arn = aws_ssoadmin_permission_set.main.arn

  principal_id   = aws_identitystore_group.main.group_id
  principal_type = "GROUP"

  target_id   = data.aws_caller_identity.current.account_id
  target_type = "AWS_ACCOUNT"
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

output "target_account_id" {
  value = resource.aws_ssoadmin_account_assignment.main.target_id
}

output "instance_arn" {
  value = resource.aws_ssoadmin_account_assignment.main.instance_arn
}

output "permission_set_arn" {
  value = resource.aws_ssoadmin_account_assignment.main.permission_set_arn
}

output "principal_id" {
  value = resource.aws_ssoadmin_account_assignment.main.principal_id
}

output "principal_type" {
  value = resource.aws_ssoadmin_account_assignment.main.principal_type
}