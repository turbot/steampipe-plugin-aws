
variable "resource_name" {
  type    = string
  default = "turbot-test-20200125-create-update"
  description = "Name of the resource used throughout the test."
}

variable "aws_profile" {
  type    = string
  default = "default"
  description = "AWS credentials profile used for the test. Default is to use the default profile."
}

variable "aws_region" {
  type    = string
  default = "us-east-1"
  description = "AWS region used for the test. Does not work with default region in config, so must be defined here."
}

variable "aws_region_alternate" {
  type    = string
  default = "us-east-2"
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

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
  tags = {
    name = var.resource_name
  }
}

resource "aws_subnet" "public" {
  vpc_id     = aws_vpc.main.id
  cidr_block = "10.0.1.0/24"
}

resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.main.id
}

resource "aws_lb" "test" {
  name               = var.resource_name
  internal           = false
  load_balancer_type = "network"
  subnets            = [aws_subnet.public.id]

  depends_on = [aws_internet_gateway.gw]

  tags = {
    name = var.resource_name
  }
}

resource "aws_vpc_endpoint_service" "named_test_resource" {
  acceptance_required        = false
  network_load_balancer_arns = [aws_lb.test.arn]
  tags = {
    Name = var.resource_name
  }
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "resource_aka" {
  value = aws_vpc_endpoint_service.named_test_resource.arn
}

output "resource_id" {
  value = aws_vpc_endpoint_service.named_test_resource.id
}

output "service_name" {
  value = aws_vpc_endpoint_service.named_test_resource.service_name
}

output "resource_name" {
  value = var.resource_name
}
