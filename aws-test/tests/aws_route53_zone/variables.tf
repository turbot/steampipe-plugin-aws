
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

# Create AWS > Route53 > Zone
resource "aws_route53_zone" "named_test_resource" {
  name          = "${var.resource_name}.com"
  comment       = "Test zone"
  force_destroy = true
  tags = {
    name = var.resource_name
  }
}

resource "aws_cloudwatch_log_group" "aws_route53_example_com" {
  name              = "/aws/route53/${aws_route53_zone.named_test_resource.name}"
  retention_in_days = 30
}

# Example CloudWatch log resource policy to allow Route53 to write logs
# to any log group under /aws/route53/*

data "aws_iam_policy_document" "route53-query-logging-policy" {
  statement {
    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    resources = ["arn:aws:logs:*:*:log-group:/aws/route53/*"]

    principals {
      identifiers = ["route53.amazonaws.com"]
      type        = "Service"
    }
  }
}

resource "aws_cloudwatch_log_resource_policy" "route53-query-logging-policy" {
  policy_document = data.aws_iam_policy_document.route53-query-logging-policy.json
  policy_name     = "route53-query-logging-policy"
}

resource "aws_route53_query_log" "example_com" {
  depends_on = [aws_cloudwatch_log_resource_policy.route53-query-logging-policy]

  cloudwatch_log_group_arn = aws_cloudwatch_log_group.aws_route53_example_com.arn
  zone_id                  = aws_route53_zone.named_test_resource.zone_id
}

resource "aws_kms_key" "example" {
  customer_master_key_spec = "ECC_NIST_P256"
  deletion_window_in_days  = 7
  key_usage                = "SIGN_VERIFY"
  policy = jsonencode({
    Statement = [
      {
        Action = [
          "kms:DescribeKey",
          "kms:GetPublicKey",
          "kms:Sign",
        ],
        Effect = "Allow"
        Principal = {
          Service = "api-service.dnssec.route53.aws.internal"
        }
        Sid = "Route 53 DNSSEC Permissions"
      },
      {
        Action = "kms:*"
        Effect = "Allow"
        Principal = {
          AWS = "*"
        }
        Resource = "*"
        Sid      = "IAM User Permissions"
      },
    ]
    Version = "2012-10-17"
  })
}

resource "aws_route53_key_signing_key" "example" {
  hosted_zone_id             = aws_route53_zone.named_test_resource.id
  key_management_service_arn = aws_kms_key.example.arn
  name                       = var.resource_name
}

resource "aws_route53_hosted_zone_dnssec" "example" {
  hosted_zone_id = aws_route53_key_signing_key.example.hosted_zone_id
}

output "resource_aka" {
  value = "arn:aws:route53:::hostedzone/${aws_route53_zone.named_test_resource.zone_id}"
}

output "zone_id" {
  value = aws_route53_zone.named_test_resource.zone_id
}

output "cloudwatch_log_group_arn" {
  value = aws_cloudwatch_log_group.aws_route53_example_com.arn
}

output "logging_configuration_id" {
  value = aws_route53_query_log.example_com.id
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

output "kms_arn" {
  value = aws_kms_key.example.arn
}

output "digest_algorithm_mnemonic" {
  value = aws_route53_key_signing_key.example.digest_algorithm_mnemonic
}

output "signing_algorithm_mnemonic" {
  value = aws_route53_key_signing_key.example.signing_algorithm_mnemonic
}
