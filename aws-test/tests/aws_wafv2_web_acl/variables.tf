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

resource "aws_wafv2_web_acl" "named_test_resource" {
  name        = var.resource_name
  description = "Example of a managed rule."
  scope       = "REGIONAL"

  default_action {
    allow {}
  }

  rule {
    name     = "rule-1"
    priority = 1

    override_action {
      count {}
    }

    statement {
      managed_rule_group_statement {
        name        = "AWSManagedRulesCommonRuleSet"
        vendor_name = "AWS"
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = false
      metric_name                = "friendly-rule-metric-${var.resource_name}"
      sampled_requests_enabled   = false
    }
  }

  tags = {
    name = var.resource_name
  }

  visibility_config {
    cloudwatch_metrics_enabled = false
    metric_name                = "friendly-metric-${var.resource_name}"
    sampled_requests_enabled   = false
  }
}

resource "aws_s3_bucket" "firehose_bucket" {
  bucket = var.resource_name
}

resource "aws_iam_role" "firehose_role" {
  name = var.resource_name

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "firehose.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_kinesis_firehose_delivery_stream" "named_test_resource" {
  name        = "aws-waf-logs-${var.resource_name}"
  destination = "extended_s3"

  extended_s3_configuration {
    role_arn   = aws_iam_role.firehose_role.arn
    bucket_arn = aws_s3_bucket.firehose_bucket.arn
  }

}

resource "aws_wafv2_web_acl_logging_configuration" "example" {
  log_destination_configs = [aws_kinesis_firehose_delivery_stream.named_test_resource.arn]
  resource_arn            = aws_wafv2_web_acl.named_test_resource.arn
  redacted_fields {
    single_header {
      name = "user-agent"
    }
  }
}

output "resource_aka" {
  value = aws_wafv2_web_acl.named_test_resource.arn
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = aws_wafv2_web_acl.named_test_resource.id
}

