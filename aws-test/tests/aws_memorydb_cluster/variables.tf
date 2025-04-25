
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
resource "aws_vpc" "named_test_resource" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "named_test_resource" {
  vpc_id            = aws_vpc.named_test_resource.id
  cidr_block        = "10.0.0.0/24"
  availability_zone = "us-east-1b"
}

resource "aws_subnet" "named_test_resource_2" {
  vpc_id            = aws_vpc.named_test_resource.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = "us-east-1c"
}

resource "aws_memorydb_subnet_group" "named_test_resource" {
  name       = "my-subnet-group"
  subnet_ids = [aws_subnet.named_test_resource.id, aws_subnet.named_test_resource_2.id]
}

resource "random_password" "example" {
  length = 16
}

resource "aws_memorydb_user" "example" {
  user_name     = var.resource_name
  access_string = "on ~* &* +@all"

  authentication_mode {
    type      = "password"
    passwords = [random_password.example.result]
  }
}

resource "aws_memorydb_acl" "named_test_resource" {
  depends_on = [ aws_memorydb_user.example ]
  name       = var.resource_name
  user_names = [var.resource_name]
}

resource "aws_memorydb_cluster" "named_test_resource" {
  depends_on = [ aws_memorydb_subnet_group.named_test_resource, aws_memorydb_acl.named_test_resource ]
  acl_name                 = var.resource_name
  name                     = var.resource_name
  node_type                = "db.t4g.small"
  num_shards               = 2
  snapshot_retention_limit = 7
  subnet_group_name        = aws_memorydb_subnet_group.named_test_resource.id
}

output "resource_aka" {
  value = aws_memorydb_cluster.named_test_resource.arn
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

