
variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125"
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

# Create AWS > Config > Configuration Recorder
resource "aws_config_configuration_recorder" "configuration_recorder" {
  name     = var.resource_name
  role_arn = aws_iam_role.r.arn
}

resource "aws_iam_role" "r" {
  name = var.resource_name

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "config.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
POLICY
}

# Create AWS > Config > Conformance pack
resource "aws_config_conformance_pack" "named_test_resource" {
  name = var.resource_name

  input_parameter {
    parameter_name  = "AccessKeysRotatedParameterMaxAccessKeyAge"
    parameter_value = "90"
  }

  template_body = <<EOT
    Parameters:
      AccessKeysRotatedParameterMaxAccessKeyAge:
        Type: String
    Resources:
      IAMPasswordPolicy:
        Properties:
          ConfigRuleName: IAMPasswordPolicy
          Source:
            Owner: AWS
            SourceIdentifier: IAM_PASSWORD_POLICY
        Type: AWS::Config::ConfigRule
EOT

  depends_on = [aws_config_configuration_recorder.configuration_recorder]
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "region_name" {
  value = data.aws_region.primary.name
}

output "aws_partition" {
  value = data.aws_partition.current.partition
}

output "resource_name" {
  value = var.resource_name
}

output "resource_aka" {
  value = aws_config_conformance_pack.named_test_resource.arn
}
