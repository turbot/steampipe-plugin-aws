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
  default     = "us-east-1"
  description = "AWS region used for the test. Does not work with default region in config, so must be defined here."
}

variable "aws_region_alternate" {
  type        = string
  default     = "us-east-2"
  description = "Alternate AWS region used for tests that require two regions (e.g. DynamoDB global tables)."
}

provider "aws" {
  // profile = var.aws_profile
  region  = var.aws_region
}

provider "aws" {
  alias   = "alternate"
  // profile = var.aws_profile
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
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "my_subnet" {
  vpc_id            = aws_vpc.my_vpc.id
  cidr_block        = "10.0.0.0/24"
  availability_zone = "${var.aws_region}a"
}

resource "aws_subnet" "my_subnet_2" {
  vpc_id            = aws_vpc.my_vpc.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = "${var.aws_region}b"
}

resource "aws_security_group" "my_security_group" {
  name_prefix = var.resource_name
  vpc_id      = aws_vpc.my_vpc.id

  ingress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_elasticache_user_group" "my_user_group" {
  engine = "valkey"
  user_group_id = var.resource_name
}

resource "aws_elasticache_serverless_cache" "named_test_resource" {
  name                 = var.resource_name
  engine               = "valkey"
  user_group_id        = aws_elasticache_user_group.my_user_group.user_group_id
  subnet_ids           = [aws_subnet.my_subnet.id, aws_subnet.my_subnet_2.id]
  security_group_ids   = [aws_security_group.my_security_group.id]
  
  cache_usage_limits {
    data_storage {
      maximum = 10
      unit    = "GB"
    }
    ecpu_per_second {
      maximum = 5000
    }
  }

  tags = {
    name = var.resource_name
  }
}

output "resource_name" {
  value = var.resource_name
}

output "resource_aka" {
  depends_on = [aws_elasticache_serverless_cache.named_test_resource]
  value      = "arn:${data.aws_partition.current.partition}:elasticache:${data.aws_region.primary.name}:${data.aws_caller_identity.current.account_id}:serverlesscache:${aws_elasticache_serverless_cache.named_test_resource.name}"
}

output "resource_id" {
  value = aws_elasticache_serverless_cache.named_test_resource.name
}

output "serverless_cache_name" {
  value = aws_elasticache_serverless_cache.named_test_resource.name
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
} 