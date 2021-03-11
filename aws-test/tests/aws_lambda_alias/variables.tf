
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

# Create AWS > Lambda > Alias
resource "local_file" "python_file" {
  filename          = "${path.cwd}/../../test.py"
  sensitive_content = "def test (event, context):\n\tprint ('This is a test for intergration testing to check creation of a lambda function')"
}

data "archive_file" "zip" {
  depends_on  = [local_file.python_file]
  type        = "zip"
  source_file = "${path.cwd}/../../test.py"
  output_path = "${path.cwd}/../../test.zip"
}

resource "aws_iam_role" "aws_lambda_function" {
  name = var.resource_name
  assume_role_policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Action" : "sts:AssumeRole",
        "Principal" : {
          "Service" : "lambda.amazonaws.com"
        },
        "Effect" : "Allow",
        "Sid" : ""
      }
    ]
  })
}

resource "aws_lambda_function" "named_test_resource" {
  function_name = var.resource_name
  role          = aws_iam_role.aws_lambda_function.arn
  handler       = "test.test"
  runtime       = "python3.7"
  filename      = "${path.cwd}/../../test.zip"
}

resource "aws_lambda_alias" "named_test_resource" {
  name             = var.resource_name
  function_name    = "arn:${data.aws_partition.current.partition}:lambda:${var.aws_region}:${data.aws_caller_identity.current.account_id}:function:${var.resource_name}"
  function_version = "$LATEST"
  description      = "Test alias."
  depends_on = [
    aws_lambda_function.named_test_resource,
  ]
}

output "resource_aka" {
  value = aws_lambda_alias.named_test_resource.arn
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "aws_partition" {
  value = data.aws_partition.current.partition
}

output "region_name" {
  value = data.aws_region.primary.name
}

output "resource_name" {
  value = var.resource_name
}
