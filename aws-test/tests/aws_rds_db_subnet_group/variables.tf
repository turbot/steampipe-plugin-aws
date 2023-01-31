
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

data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_vpc" "aws_db_subnet_group" {
  cidr_block = "10.84.0.0/24"
}

resource "aws_subnet" "aws_db_subnet_group_a" {
  vpc_id            = aws_vpc.aws_db_subnet_group.id
  cidr_block        = "10.84.0.0/25"
  availability_zone = data.aws_availability_zones.available.names[0]
}

resource "aws_subnet" "aws_db_subnet_group_b" {
  vpc_id            = aws_vpc.aws_db_subnet_group.id
  cidr_block        = "10.84.0.128/25"
  availability_zone = data.aws_availability_zones.available.names[1]
}

resource "aws_db_subnet_group" "named_test_resource" {
  name        = var.resource_name
  description = "Terraform created resource to verify table"
  subnet_ids  = [aws_subnet.aws_db_subnet_group_a.id, aws_subnet.aws_db_subnet_group_b.id]
  tags = {
    name = var.resource_name
  }
}

output "vpc_id" {
  value = aws_vpc.aws_db_subnet_group.id
}

output "resource_name" {
  value = var.resource_name
}

output "resource_aka" {
  value = aws_db_subnet_group.named_test_resource.arn
}
