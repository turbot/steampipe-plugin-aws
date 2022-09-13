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

resource "aws_amplify_app" "named_test_resource" {
  name                        = var.resource_name
  description                 = "Description field"
  enable_auto_branch_creation = true
  enable_basic_auth           = true
  enable_branch_auto_deletion = true
  basic_auth_credentials      = base64encode("username1:password1")

  # The default build_spec added by the Amplify Console for React.
  build_spec = <<-EOT
version: 1
backend:
  phases:
    build:
      commands:
        - "build command"
EOT

  platform = "WEB"
  tags = {
    Purpose     = "Testing"
    Application = "Steampipe"
  }

  custom_rule {
    source = "/<*>"
    status = "404"
    target = "/index.html"
  }

  environment_variables = {
    TOOL    = "steampipe"
    PURPOSE = "testing"
  }
}

resource "time_sleep" "wait_120_seconds" {
  depends_on = [aws_amplify_app.named_test_resource]

  create_duration = "120s"
}

output "id" {
  value = aws_amplify_app.named_test_resource.id
}

output "resource_aka" {
  value = aws_amplify_app.named_test_resource.arn
}

output "resource_name" {
  value = var.resource_name
}
