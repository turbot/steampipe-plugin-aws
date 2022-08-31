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
  default     = "us-east-1"
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

resource "aws_iam_group" "named_test_resource" {
  name = var.resource_name
}

resource "aws_iam_group_policy" "group_inline_policy" {
  name   = var.resource_name
  group  = aws_iam_group.named_test_resource.name
  policy = "{\n  \"Version\": \"2012-10-17\",\n  \"Statement\": [\n    {\n      \"Action\": [\n        \"ec2:Describe*\"\n      ],\n      \"Effect\": \"Allow\",\n      \"Resource\": \"*\"\n    }\n  ]\n}\n"
}

# resource "aws_iam_user" "named_test_user" {
#   name = var.resource_name
#   tags = {
#     name = var.resource_name
#   }
# }

# resource "aws_iam_group_membership" "team" {
#   name = "tf-testing-group-membership"

#   users = [
#     aws_iam_user.named_test_user.name,
#   ]

#   group = aws_iam_group.named_test_resource.name
# }

resource "aws_iam_policy" "policy" {
  name        = var.resource_name
  description = "A test policy"
  policy      = "{\n  \"Version\": \"2012-10-17\",\n  \"Statement\": [\n    {\n      \"Action\": [\n        \"ec2:Describe*\"\n      ],\n      \"Effect\": \"Allow\",\n      \"Resource\": \"*\"\n    }\n  ]\n}\n"
}

resource "aws_iam_group_policy_attachment" "test_attach" {
  group      = aws_iam_group.named_test_resource.name
  policy_arn = aws_iam_policy.policy.arn
}


output "attached_policy_arn" {
  value = aws_iam_group_policy_attachment.test_attach.policy_arn
}

output "resource_aka" {
  value = aws_iam_group.named_test_resource.arn
}

output "resource_name" {
  value = var.resource_name
}

output "group_id" {
  value = aws_iam_group.named_test_resource.unique_id
}

output "aws_partition" {
  value = data.aws_partition.current.partition
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}
