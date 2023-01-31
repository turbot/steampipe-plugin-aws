
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

resource "aws_wafv2_rule_group" "named_test_resource_regional" {
  name        = "${var.resource_name}_regional"
  description = "A regional rule group for testing."
  scope       = "REGIONAL"
  capacity    = 2

  rule {
    name     = "rule-1"
    priority = 1
    action {
      allow {}
    }
    statement {
      geo_match_statement {
        country_codes = ["US", "IN"]
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = false
      metric_name                = "${var.resource_name}_regional"
      sampled_requests_enabled   = false
    }
  }

  visibility_config {
    cloudwatch_metrics_enabled = false
    metric_name                = "${var.resource_name}_regional"
    sampled_requests_enabled   = false
  }

  tags = {
    foo   = "bar"
    scope = "regional"
  }
}

resource "aws_wafv2_rule_group" "named_test_resource_global" {
  name        = "${var.resource_name}_global"
  description = "A global rule group for testing."
  scope       = "CLOUDFRONT"
  capacity    = 2

  rule {
    name     = "rule-1"
    priority = 1
    action {
      allow {}
    }
    statement {
      geo_match_statement {
        country_codes = ["US", "IN"]
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = false
      metric_name                = "${var.resource_name}_global"
      sampled_requests_enabled   = false
    }
  }

  visibility_config {
    cloudwatch_metrics_enabled = false
    metric_name                = "${var.resource_name}_global"
    sampled_requests_enabled   = false
  }

  tags = {
    foo1  = "bar1"
    scope = "global"
  }
}

output "resource_aka_regional" {
  value = aws_wafv2_rule_group.named_test_resource_regional.arn
}

output "resource_aka_global" {
  value = aws_wafv2_rule_group.named_test_resource_global.arn
}

output "resource_id_regional" {
  value = aws_wafv2_rule_group.named_test_resource_regional.id
}

output "resource_id_global" {
  value = aws_wafv2_rule_group.named_test_resource_global.id
}

output "resource_name_regional" {
  value = "${var.resource_name}_regional"
}

output "resource_name_global" {
  value = "${var.resource_name}_global"
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
