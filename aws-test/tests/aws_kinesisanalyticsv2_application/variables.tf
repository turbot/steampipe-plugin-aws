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

resource "aws_iam_role" "named_test_resource" {
  name               = var.resource_name
  assume_role_policy = "{\n  \"Version\": \"2012-10-17\",\n  \"Statement\": [\n    {\n      \"Action\": \"sts:AssumeRole\",\n      \"Principal\": {\n        \"Service\": \"kinesisanalytics.amazonaws.com\"\n      },\n      \"Effect\": \"Allow\",\n      \"Sid\": \"test\"\n    }\n  ]\n}\n"
}

resource "aws_kinesisanalyticsv2_application" "named_test_resource" {
  name                   = var.resource_name
  runtime_environment    = "SQL-1_0"
  service_execution_role = aws_iam_role.named_test_resource.arn
  tags = {
    Name = var.resource_name
  }
}

output "resource_aka" {
  value = aws_kinesisanalyticsv2_application.named_test_resource.arn
}

output "resource_id" {
  value = aws_kinesisanalyticsv2_application.named_test_resource.id
}

output "resource_creation_time" {
  value = aws_kinesisanalyticsv2_application.named_test_resource.create_timestamp
}

output "version" {
  value = aws_kinesisanalyticsv2_application.named_test_resource.version_id
}

output "resource_name" {
  value = var.resource_name
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "aws_partition" {
  value = data.aws_partition.current.partition
}

output "aws_region" {
  value = var.aws_region
}
