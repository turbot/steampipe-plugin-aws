
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

variable "auth_mode" {
  type        = string
  default     = "IAM"
  description = "The mode of authentication that members use to access the domain."
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

resource "aws_iam_role" "my_role" {
  name = var.resource_name
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "sagemaker.amazonaws.com"
        }
      },
    ]
  })
}

resource "aws_vpc" "named_test_resource" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "named_test_resource" {
  vpc_id     = aws_vpc.named_test_resource.id
  cidr_block = "10.0.0.0/24"
}

resource "aws_sagemaker_domain" "named_test_resource" {
  domain_name = var.resource_name
  auth_mode   = var.auth_mode
  vpc_id      = aws_vpc.named_test_resource.id
  subnet_ids  = [aws_subnet.named_test_resource.id]

  default_user_settings {
    execution_role = aws_iam_role.my_role.arn
  }
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

output "resource_aka" {
  value = aws_sagemaker_domain.named_test_resource.arn
}

output "resource_id" {
  value = aws_sagemaker_domain.named_test_resource.id
}

output "resource_name" {
  value = var.resource_name
}
