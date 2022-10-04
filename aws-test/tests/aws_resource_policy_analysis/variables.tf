variable "resource_name" {
  type        = string
  default     = "turbot-test-20221004"
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

resource "aws_iam_role" "test_role" {
  name = var.resource_name

  # Terraform's "jsonencode" function converts a
  # Terraform expression result to valid JSON syntax.
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = "PublicService"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      },
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = "RestrictedService"
        Principal = {
          Service = "cloudwatch.amazonaws.com"
        }
        Condition = {
          StringEquals = {
            "aws:SourceOwner" = "012345678901"
          }
        }
      },
    ]
  })

  tags = {
    tag-key = "integration-test"
  }
}

output "resource_aka" {
  value = aws_iam_role.test_role.arn
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = aws_iam_role.test_role.unique_id
}
