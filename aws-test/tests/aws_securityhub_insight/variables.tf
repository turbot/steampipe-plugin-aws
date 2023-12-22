
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
  default     = "us-west-2"
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

resource "aws_securityhub_account" "named_test_resource" {}

// Get the securityhub insights
resource "null_resource" "list_insights" {
  provisioner "local-exec" {
    command = "aws securityhub create-insight --region ${data.aws_region.primary.name} --filters '{\"ResourceType\": [{\"Comparison\": \"EQUALS\", \"Value\": \"AwsIamRole\"}], \"SeverityLabel\": [{\"Comparison\": \"EQUALS\", \"Value\": \"CRITICAL\"}]}' --group-by-attribute \"ResourceId\" --name ${var.resource_name} --profile ${var.aws_profile}"
  }
  provisioner "local-exec" {
    command = "aws securityhub get-insights --output json --profile ${var.aws_profile} --region ${data.aws_region.primary.name} > ${path.cwd}/list-insights.json"
  }
}

data "local_file" "insights" {
  depends_on = [null_resource.list_insights]
  filename   = "${path.cwd}/list-insights.json"
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "insight_arn" {
  value = jsondecode(data.local_file.insights.content).Insights[0].InsightArn
}

output "group_by_att" {
  value = jsondecode(data.local_file.insights.content).Insights[0].GroupByAttribute
}

output "name" {
  value = jsondecode(data.local_file.insights.content).Insights[0].Name
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
