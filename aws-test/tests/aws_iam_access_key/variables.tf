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

resource "aws_iam_user" "named_test_resource" {
  name = var.resource_name
  tags = {
    name = var.resource_name
  }
}

resource "aws_iam_access_key" "named_test_resource" {
  user = aws_iam_user.named_test_resource.name
}

locals {
  resource_aka = "arn:${data.aws_partition.current.partition}:iam::${data.aws_caller_identity.current.account_id}:user/${aws_iam_user.named_test_resource.name}/accesskey/${aws_iam_access_key.named_test_resource.id}"
}

output "resource_aka" {
  depends_on = [aws_iam_access_key.named_test_resource]
  value      = local.resource_aka
}

output "resource_id" {
  value = aws_iam_access_key.named_test_resource.id
}

output "user_name" {
  value = aws_iam_user.named_test_resource.name
}

output "aws_partition" {
  value = data.aws_partition.current.partition
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}
