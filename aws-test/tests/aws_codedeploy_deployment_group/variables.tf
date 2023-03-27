variable "resource_name" {
  type        = string
  default     = "turbot-test-20200126"
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
  description = "Alternate AWS region used for tests that require two regions."
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

data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["codedeploy.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "named_test_resource" {
  name               = var.resource_name
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

resource "aws_iam_role_policy_attachment" "AWSCodeDeployRole" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSCodeDeployRole"
  role       = aws_iam_role.named_test_resource.name
}

resource "aws_codedeploy_app" "named_test_resource" {
  name = var.resource_name
}

resource "aws_codedeploy_deployment_group" "named_test_resource" {
  app_name              = aws_codedeploy_app.named_test_resource.name
  deployment_group_name = var.resource_name
  service_role_arn      = aws_iam_role.named_test_resource.arn
  tags = {
    name = var.resource_name
  }
}

output "resource_aka" {
  value = aws_codedeploy_deployment_group.named_test_resource.arn
}

output "app_name" {
  value = aws_codedeploy_deployment_group.named_test_resource.app_name
}

output "resource_id" {
  value = aws_codedeploy_deployment_group.named_test_resource.deployment_group_id
}

output "resource_name" {
  value = aws_codedeploy_deployment_group.named_test_resource.deployment_group_name
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}