
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
  template = "arn:$${partition}:ec2:$${region}:$${account_id}:route-table/${aws_route_table.named_test_resource.id}"
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

resource "aws_subnet" "named_test_resource" {
  vpc_id     = aws_vpc.main.id
  cidr_block = "10.0.1.0/24"
  tags = {
    Name = var.resource_name
  }
}

resource "aws_route_table" "named_test_resource" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = var.resource_name
  }
}

resource "aws_route_table_association" "named_test_resource" {
  subnet_id      = aws_subnet.named_test_resource.id
  route_table_id = aws_route_table.named_test_resource.id
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "vpc_id" {
  value = aws_vpc.main.id
}

output "subnet_id" {
  value = aws_subnet.named_test_resource.id
}

output "association_id" {
  value = aws_route_table_association.named_test_resource.id
}

output "resource_id" {
  value = aws_route_table.named_test_resource.id
}

output "resource_name" {
  value = var.resource_name
}

output "resource_aka" {
  depends_on = [aws_route_table.named_test_resource]
  value      = data.template_file.resource_aka.rendered
}
