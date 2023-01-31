

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

# Create AWS > SQS > QUEUE
resource "aws_sqs_queue" "named_test_resource" {
  name = var.resource_name
  tags = {
    name = var.resource_name
  }
}

# Create AWS > SNS > Topic
resource "aws_sns_topic" "named_test_resource" {
  name         = var.resource_name
  display_name = var.resource_name
}

resource "aws_sqs_queue_policy" "test" {
  queue_url = aws_sqs_queue.named_test_resource.id

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Id": "sqspolicy",
  "Statement": [
    {
      "Sid": "First",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "sqs:SendMessage",
      "Resource": "${aws_sqs_queue.named_test_resource.arn}",
      "Condition": {
        "ArnEquals": {
          "aws:SourceArn": "${aws_sns_topic.named_test_resource.arn}"
        }
      }
    }
  ]
}
POLICY
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "aws_region" {
  value = data.aws_region.primary.name
}


output "resource_aka" {
  value = aws_sqs_queue.named_test_resource.arn
}

output "sns_topic_arn" {
  value = aws_sns_topic.named_test_resource.arn
}

output "resource_name" {
  value = var.resource_name
}

output "queue_url" {
  value = aws_sqs_queue.named_test_resource.id
}
