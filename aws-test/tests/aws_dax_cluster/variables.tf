
variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125"
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

resource "aws_iam_role" "my_role" {
  name               = var.resource_name
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "dax.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_vpc" "my_vpc" {
  cidr_block = "10.1.0.0/16"
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

resource "aws_dax_subnet_group" "my_subnet_group" {
  name       = var.resource_name
  subnet_ids = [aws_subnet.my_subnet1.id, aws_subnet.my_subnet2.id]
}

resource "aws_dax_cluster" "named_test_resource" {
  cluster_name       = var.resource_name
  description        = "A test cluster"
  iam_role_arn       = aws_iam_role.my_role.arn
  node_type          = "dax.r4.large"
  replication_factor = 1
  subnet_group_name  = aws_dax_subnet_group.my_subnet_group.name
  tags = {
    name = var.resource_name
  }
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "aws_region" {
  value = data.aws_region.primary.name
}

output "aws_partition" {
  value = data.aws_partition.current.partition
}

output "resource_aka" {
  value = aws_dax_cluster.named_test_resource.arn
}

output "resource_id" {
  value = aws_dax_cluster.named_test_resource.id
}

output "resource_name" {
  value = var.resource_name
}
