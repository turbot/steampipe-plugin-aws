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


resource "aws_iam_openid_connect_provider" "named_test_resource" {
  url = "https://accounts.google.com"
  client_id_list = [
    "266362248691-342342xasdasdasda-apps.googleusercontent.com",
  ]
  thumbprint_list = ["cf23df2207d99a74fbe169e3eba035e633b65d94", "6938fd4d98bab03faadb97b34396831e3780aea1"]
  tags = {
    name = var.resource_name
  }
}


output "resource_aka" {
  value = aws_iam_openid_connect_provider.named_test_resource.arn
}

output "resource_name" {
  value = var.resource_name
}

output "aws_partition" {
  value = data.aws_partition.current.partition
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "aws_region" {
  value = "global"
}
