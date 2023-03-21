
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

resource "aws_api_gateway_domain_name" "example" {
  certificate_arn = aws_acm_certificate_validation.example.certificate_arn
  domain_name     = "api.example.com"
}

# Example DNS record using Route53.
# Route53 is not specifically required; any DNS host can be used.
resource "aws_route53_record" "example" {
  name    = aws_api_gateway_domain_name.example.domain_name
  type    = "A"
  zone_id = aws_route53_zone.example.id

  alias {
    evaluate_target_health = true
    name                   = aws_api_gateway_domain_name.example.cloudfront_domain_name
    zone_id                = aws_api_gateway_domain_name.example.cloudfront_zone_id
  }
}

output "resource_aka" {
  value = aws_api_gateway_domain_name.named_test_resource.arn
}

output "resource_id" {
  value = aws_api_gateway_domain_name.named_test_resource.id
}

output "resource_name" {
  value = "${var.resource_name}.com"
}
