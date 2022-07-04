variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update"
  description = "Name of the resource used throughout the test."
}

variable "aws_profile" {
  type        = string
  default     = "integration-tests"
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

resource "aws_iam_user" "named_test_resource" {
  name = var.resource_name
  tags = {
    name = var.resource_name
  }
}

resource "aws_iam_access_key" "named_test_resource" {
  user = aws_iam_user.named_test_resource.name
}

# Build the resource AKA dynamically based on various information about the
# account, region, zones, etc. This is important for resources where the AKA
# is not made available through the terraform provider resource definition.
# We use a terraform template to make it easy to specify substitutions.
# Be careful to escape the $ in the resourceAka definition in YAML.
# e.g. from YAML - resourceAka: "arn:$${partition}:events:$${region}:$${account_id}:rule/$${resource_name}/target/$${resource_name}"
data "template_file" "resource_aka" {
  template = "arn:$${partition}:iam::$${account_id}:user/$${resource_name}/accesskey/${aws_iam_access_key.named_test_resource.id}"
  vars = {
    resource_name    = aws_iam_user.named_test_resource.name
    partition        = data.aws_partition.current.partition
    account_id       = data.aws_caller_identity.current.account_id
    region           = data.aws_region.primary.name
    alternate_region = data.aws_region.alternate.name
  }
}

output "resource_aka" {
  depends_on = [aws_iam_access_key.named_test_resource]
  value      = data.template_file.resource_aka.rendered
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
