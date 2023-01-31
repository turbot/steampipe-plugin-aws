
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

data "template_file" "resource_aka" {
  template = "arn:$${partition}:ec2:$${region}:$${account_id}:natgateway/${aws_nat_gateway.named_test_resource.id}"
  vars = {
    partition  = data.aws_partition.current.partition
    account_id = data.aws_caller_identity.current.account_id
    region     = data.aws_region.primary.name
  }
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
  tags = {
    name = var.resource_name
  }
}

resource "aws_subnet" "example" {
  vpc_id     = aws_vpc.main.id
  cidr_block = "10.0.1.0/24"
  tags = {
    name = var.resource_name
  }
}

resource "aws_eip" "example" {
  tags = {
    Name = var.resource_name
  }
  vpc = true
}

resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.main.id
}

resource "aws_nat_gateway" "named_test_resource" {
  allocation_id = aws_eip.example.id
  subnet_id     = aws_subnet.example.id
  tags = {
    Name = var.resource_name
  }
  depends_on = [aws_internet_gateway.gw]
}

output "vpc_id" {
  value = aws_vpc.main.id
}

output "subnet_id" {
  value = aws_subnet.example.id
}

output "allocation_id" {
  value = aws_nat_gateway.named_test_resource.allocation_id
}

output "network_interface_id" {
  value = aws_nat_gateway.named_test_resource.network_interface_id
}

output "private_ip" {
  value = aws_nat_gateway.named_test_resource.private_ip
}

output "public_ip" {
  value = aws_nat_gateway.named_test_resource.public_ip
}

output "resource_id" {
  value = aws_nat_gateway.named_test_resource.id
}

output "resource_name" {
  value = var.resource_name
}

output "resource_aka" {
  depends_on = [aws_nat_gateway.named_test_resource]
  value      = data.template_file.resource_aka.rendered
}

