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

resource "aws_glue_catalog_database" "named_test_resource" {
  depends_on = [aws_iam_role.named_test_resource]
  name       = var.resource_name
}

resource "aws_iam_role" "named_test_resource" {
  name = var.resource_name
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "glue.amazonaws.com"
        }
      },
    ]
  })
}

resource "aws_dynamodb_table" "named_test_resource" {
  depends_on = [aws_glue_catalog_database.named_test_resource]
  name       = var.resource_name
  tags = {
    name = var.resource_name
  }
  hash_key = "userId"

  attribute {
    name = "userId"
    type = "S"
  }
  write_capacity = 20
  read_capacity  = 20
}

resource "aws_glue_crawler" "named_test_resource" {
  depends_on    = [aws_dynamodb_table.named_test_resource]
  database_name = aws_glue_catalog_database.named_test_resource.name
  name          = var.resource_name
  role          = aws_iam_role.named_test_resource.arn
  description   = "integration testing"

  dynamodb_target {
    path = var.resource_name
  }
}

output "resource_aka" {
  value = aws_glue_crawler.named_test_resource.arn
}

output "aws_region" {
  value = data.aws_region.primary.name
}

output "aws_partition" {
  value = data.aws_partition.current.partition
}

output "resource_name" {
  value = var.resource_name
}

output "aws_account" {
  value = data.aws_caller_identity.current.account_id
}
