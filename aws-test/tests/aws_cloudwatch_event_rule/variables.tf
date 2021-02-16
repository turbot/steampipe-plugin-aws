
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

resource "aws_cloudwatch_event_rule" "named_test_resource" {
  name        = var.resource_name
  description = "Capture each AWS Console Sign In"
  tags = {
    name = var.resource_name
  }
  event_pattern = <<EOF
{
  "detail-type": [
    "AWS Console Sign In via CloudTrail"
  ]
}
EOF
}

resource "aws_cloudwatch_event_target" "named_test_resource" {
  rule      = aws_cloudwatch_event_rule.named_test_resource.name
  target_id = "SendToSNS"
  arn       = aws_sns_topic.named_test_resource.arn
}

resource "aws_sns_topic" "named_test_resource" {
  name = var.resource_name
}

resource "aws_sns_topic_policy" "named_test_resource" {
  arn    = aws_sns_topic.named_test_resource.arn
  policy = data.aws_iam_policy_document.named_test_resource.json
}

data "aws_iam_policy_document" "named_test_resource" {
  statement {
    effect  = "Allow"
    actions = ["SNS:Publish"]

    principals {
      type        = "Service"
      identifiers = ["events.amazonaws.com"]
    }

    resources = [aws_sns_topic.named_test_resource.arn]
  }
}
output "event_pattern" {
  value = aws_cloudwatch_event_rule.named_test_resource.event_pattern
}

output "resource_aka" {
  value = aws_cloudwatch_event_rule.named_test_resource.arn
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
