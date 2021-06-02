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

resource "aws_cloudfront_origin_access_identity" "named_test_resource" {
  comment = var.resource_name
}

output "resource_aka" {
  value = "arn:${data.aws_partition.current.partition}:cloudfront::${data.aws_caller_identity.current.account_id}:origin-access-identity/${aws_cloudfront_origin_access_identity.named_test_resource.id}"
}

output "resource_id" {
  value = aws_cloudfront_origin_access_identity.named_test_resource.id
}

output "s3_canonical_user_id" {
  value = aws_cloudfront_origin_access_identity.named_test_resource.s3_canonical_user_id
}

output "caller_reference" {
  value = aws_cloudfront_origin_access_identity.named_test_resource.caller_reference
}

output "etag" {
  value = aws_cloudfront_origin_access_identity.named_test_resource.etag
}

output "resource_name" {
  value = var.resource_name
}