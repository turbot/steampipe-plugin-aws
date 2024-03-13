
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

resource "aws_networkfirewall_firewall_policy" "named_test_resource" {
  name        = var.resource_name
  description = var.resource_name

  firewall_policy {
    stateless_default_actions          = ["aws:pass"]
    stateless_fragment_default_actions = ["aws:drop"]
  }

  tags = {
    name = var.resource_name
  }
}

resource "aws_vpc" "named_test_resource" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "named_test_resource" {
  vpc_id            = aws_vpc.named_test_resource.id
  availability_zone = "${var.aws_region}a"
  cidr_block        = "10.0.1.0/24"
}

resource "aws_networkfirewall_firewall" "named_test_resource" {
  name                = var.resource_name
  firewall_policy_arn = aws_networkfirewall_firewall_policy.named_test_resource.arn
  vpc_id              = aws_vpc.named_test_resource.id
  subnet_mapping {
    subnet_id = aws_subnet.named_test_resource.id
  }

  tags = {
    name = var.resource_name
  }
}


output "resource_id" {
  value = aws_networkfirewall_firewall.named_test_resource.id
}

output "resource_aka" {
  value = aws_networkfirewall_firewall.named_test_resource.arn
}

output "resource_tags" {
  value = aws_networkfirewall_firewall.named_test_resource.tags_all
}

output "resource_vpc_id" {
  value = aws_vpc.named_test_resource.id
}

output "resource_policy_arn" {
  value = aws_networkfirewall_firewall_policy.named_test_resource.arn
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
