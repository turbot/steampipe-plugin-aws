
variable "resource_name" {
  type        = string
  default     = "turbottest-20200125"
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

# Create two subnets dynamically across different AZs
data "aws_availability_zones" "available" {}

# Create VPC
resource "aws_vpc" "test_vpc" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "${var.resource_name}-vpc"
  }
}

resource "aws_subnet" "test_subnets" {
  count                   = 2
  vpc_id                  = aws_vpc.test_vpc.id
  cidr_block              = cidrsubnet(aws_vpc.test_vpc.cidr_block, 8, count.index)
  availability_zone       = data.aws_availability_zones.available.names[count.index]

  tags = {
    Name = "${var.resource_name}-subnet-${count.index}"
  }
}

# Create DMS Replication Subnet Group
resource "aws_dms_replication_subnet_group" "test_subnet_group" {
  replication_subnet_group_id          = "${var.resource_name}-subnet-group"
  replication_subnet_group_description = "DMS replication subnet group for testing"
  subnet_ids                           = [for s in aws_subnet.test_subnets : s.id]

  tags = {
    Name = "${var.resource_name}-dms-subnet-group"
  }
}

resource "aws_dms_replication_instance" "named_test_resource" {
  allocated_storage            = 5
  apply_immediately            = true
  auto_minor_version_upgrade   = true
  availability_zone            = "${var.aws_region}a"
  multi_az                     = false
  preferred_maintenance_window = "sun:10:30-sun:14:30"
  publicly_accessible          = false
  replication_instance_class   = "dms.t3.small"
  replication_instance_id      = var.resource_name
  replication_subnet_group_id  = aws_dms_replication_subnet_group.test_subnet_group.id
  tags = {
    foo = "bar"
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
  value = aws_dms_replication_instance.named_test_resource.replication_instance_arn
}

output "resource_name" {
  value = var.resource_name
}
