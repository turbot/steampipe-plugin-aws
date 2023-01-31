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

resource "aws_vpc" "my_vpc" {
  cidr_block = "10.1.0.0/16"
}
resource "aws_security_group" "my_security_group_1" {
  name        = "create-endpoint-001"
  description = "used for creating R53 endpoint"
  vpc_id      = aws_vpc.my_vpc.id
}
resource "aws_security_group" "my_security_group_2" {
  name        = "create-endpoint-002"
  description = "used for creating R53 endpoint"
  vpc_id      = aws_vpc.my_vpc.id
}
resource "aws_security_group" "my_security_group_3" {
  name        = "create-endpoint-003"
  description = "used for creating R53 endpoint"
  vpc_id      = aws_vpc.my_vpc.id
}
resource "aws_subnet" "my_subnet1" {
  cidr_block        = "10.1.1.0/24"
  availability_zone = "${var.aws_region}a"
  vpc_id            = aws_vpc.my_vpc.id
}

resource "aws_subnet" "my_subnet2" {
  cidr_block        = "10.1.2.0/24"
  availability_zone = "${var.aws_region}b"
  vpc_id            = aws_vpc.my_vpc.id
}

resource "aws_subnet" "my_subnet3" {
  cidr_block        = "10.1.3.0/24"
  availability_zone = "${var.aws_region}c"
  vpc_id            = aws_vpc.my_vpc.id
}

resource "aws_route53_resolver_endpoint" "named_test_resource" {
  name      = var.resource_name
  direction = "INBOUND"

  tags = {
    name = var.resource_name
  }
  security_group_ids = [
    aws_security_group.my_security_group_1.id,
    aws_security_group.my_security_group_2.id,
    aws_security_group.my_security_group_3.id,
  ]

  ip_address {
    subnet_id = aws_subnet.my_subnet1.id
  }

  ip_address {
    subnet_id = aws_subnet.my_subnet3.id
  }

  ip_address {
    subnet_id = aws_subnet.my_subnet2.id
    ip        = "10.1.2.4"
  }
}

output "resource_id" {
  value = aws_route53_resolver_endpoint.named_test_resource.id
}

output "resource_aka" {
  value = aws_route53_resolver_endpoint.named_test_resource.arn
}

output "resource_direction" {
  value = aws_route53_resolver_endpoint.named_test_resource.direction
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
