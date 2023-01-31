
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

resource "aws_route53_traffic_policy" "named_test_resource" {
  name     = var.resource_name
  comment  = "${var.resource_name} comment"
  document = <<EOF
    {
      "AWSPolicyFormatVersion": "2015-10-01",
      "RecordType": "A",
      "Endpoints": {
        "endpoint-start-NkPh": {
          "Type": "value",
          "Value": "10.0.0.2"
        }
      },
      "StartEndpoint": "endpoint-start-NkPh"
    }
    EOF
}

output "id" {
  value = aws_route53_traffic_policy.named_test_resource.id
}

output "type" {
  value = aws_route53_traffic_policy.named_test_resource.type
}

output "version" {
  value = aws_route53_traffic_policy.named_test_resource.version
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "aws_region" {
  value = data.aws_region.primary.name
}

output "resource_name" {
  value = var.resource_name
}
