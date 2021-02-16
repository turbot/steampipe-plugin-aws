
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

# Create AWS > Elastic IP
resource "aws_eip" "named_test_resource" {
  vpc = true
}

# Create AWS > Route53 > Zone
resource "aws_route53_zone" "named_test_resource" {
  name          = "${var.resource_name}.com"
  force_destroy = true
}

# Create AWS > Route53 > Zone > Records
resource "aws_route53_record" "named_test_resource" {
  zone_id = aws_route53_zone.named_test_resource.zone_id
  name    = "www.${var.resource_name}.com"
  type    = "A"
  ttl     = "300"

  weighted_routing_policy {
    weight = 90
  }

  set_identifier = "live"
  records        = [aws_eip.named_test_resource.public_ip]
}

output "resource_aka" {
  value = "arn:aws:route53:::hostedzone/${aws_route53_zone.named_test_resource.zone_id}/recordset/www.${var.resource_name}.com./A/live"
}

output "zone_id" {
  value = aws_route53_zone.named_test_resource.zone_id
}

output "public_ip" {
  value = aws_eip.named_test_resource.public_ip
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
