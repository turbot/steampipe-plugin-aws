
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

resource "aws_vpc" "my_vpc" {
  cidr_block = "10.1.0.0/16"
  tags = {
    name = var.resource_name
  }
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.my_vpc.id
  tags = {
    name = var.resource_name
  }
}

resource "aws_subnet" "my_subnet1" {
  cidr_block        = "10.1.1.0/24"
  availability_zone = "${var.aws_region}a"
  vpc_id            = aws_vpc.my_vpc.id
  depends_on        = [aws_internet_gateway.igw]
  tags = {
    name = var.resource_name
  }
}

resource "aws_subnet" "my_subnet2" {
  cidr_block        = "10.1.2.0/24"
  availability_zone = "${var.aws_region}b"
  vpc_id            = aws_vpc.my_vpc.id
  depends_on        = [aws_internet_gateway.igw]
  tags = {
    name = var.resource_name
  }
}

resource "aws_db_subnet_group" "my_subnetgroup" {
  name = var.resource_name
  subnet_ids = [
    aws_subnet.my_subnet1.id,
    aws_subnet.my_subnet2.id
  ]
}

resource "aws_rds_cluster" "test_cluster" {
  cluster_identifier   = var.resource_name
  master_username      = "turbottest"
  master_password      = "test123Q"
  skip_final_snapshot  = true
  db_subnet_group_name = aws_db_subnet_group.my_subnetgroup.name
  tags = {
    name = var.resource_name
  }
}

resource "aws_db_cluster_snapshot" "named_test_resource" {
  db_cluster_snapshot_identifier = var.resource_name
  db_cluster_identifier          = aws_rds_cluster.test_cluster.id
  tags = {
    name = var.resource_name
  }
}

output "vpc_id" {
  value = aws_vpc.my_vpc.id
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "aws_partition" {
  value = data.aws_partition.current.partition
}

output "region_name" {
  value = data.aws_region.primary.name
}

output "resource_name" {
  value = var.resource_name
}

output "resource_aka" {
  value = aws_db_cluster_snapshot.named_test_resource.db_cluster_snapshot_arn
}
