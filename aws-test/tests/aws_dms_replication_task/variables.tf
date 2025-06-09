
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

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
  tags = {
    name = var.resource_name
  }
}

resource "aws_subnet" "sub1" {
  vpc_id            = aws_vpc.main.id
  availability_zone = "us-east-1a"
  cidr_block        = "10.0.0.0/17"
  tags = {
    Name = var.resource_name
  }
}
resource "aws_subnet" "sub2" {
  vpc_id            = aws_vpc.main.id
  availability_zone = "us-east-1b"
  cidr_block        = "10.0.128.0/17"
  tags = {
    Name = var.resource_name
  }
}

# Create a new replication subnet group
resource "aws_dms_replication_subnet_group" "named_test_resource" {
  replication_subnet_group_description = "Example replication subnet group"
  replication_subnet_group_id          = var.resource_name

  subnet_ids = [
    aws_subnet.sub1.id,
    aws_subnet.sub2.id,
  ]
}

resource "aws_dms_replication_instance" "named_test_resource" {
  depends_on                   = [aws_dms_replication_subnet_group.named_test_resource]
  allocated_storage            = 5
  apply_immediately            = true
  auto_minor_version_upgrade   = true
  availability_zone            = "${var.aws_region}a"
  multi_az                     = false
  preferred_maintenance_window = "sun:10:30-sun:14:30"
  publicly_accessible          = false
  replication_instance_class   = "dms.t3.small"
  replication_instance_id      = var.resource_name
  replication_subnet_group_id  = aws_dms_replication_subnet_group.named_test_resource.id
  tags = {
    foo = "bar"
  }
}

resource "aws_dms_endpoint" "named_test_resource_target" {
  depends_on                  = [aws_dms_replication_instance.named_test_resource]
  database_name               = "test"
  endpoint_id                 = var.resource_name
  endpoint_type               = "target"
  engine_name                 = "aurora"
  extra_connection_attributes = ""
  password                    = "test"
  port                        = 3306
  server_name                 = "test"
  ssl_mode                    = "none"
  username                    = "test"
}
resource "aws_dms_endpoint" "named_test_resource_source" {
  depends_on                  = [aws_dms_replication_instance.named_test_resource]
  database_name               = "test"
  endpoint_id                 = "${var.resource_name}source"
  endpoint_type               = "source"
  engine_name                 = "aurora"
  extra_connection_attributes = ""
  password                    = "test"
  port                        = 3305
  server_name                 = "test"
  ssl_mode                    = "none"
  username                    = "test"
}


# Create a new replication task
resource "aws_dms_replication_task" "named_test_resource" {
  depends_on               = [aws_dms_replication_instance.named_test_resource, aws_dms_endpoint.named_test_resource_source, aws_dms_endpoint.named_test_resource_target]
  cdc_start_time           = "1993-05-21T05:50:00Z"
  migration_type           = "full-load"
  replication_instance_arn = aws_dms_replication_instance.named_test_resource.replication_instance_arn
  replication_task_id      = var.resource_name
  source_endpoint_arn      = aws_dms_endpoint.named_test_resource_source.endpoint_arn
  table_mappings           = "{\"rules\":[{\"rule-type\":\"selection\",\"rule-id\":\"1\",\"rule-name\":\"1\",\"object-locator\":{\"schema-name\":\"%\",\"table-name\":\"%\"},\"rule-action\":\"include\"}]}"

  tags = {
    Name = "test"
  }

  target_endpoint_arn = aws_dms_endpoint.named_test_resource_target.endpoint_arn
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
  value = aws_dms_replication_task.named_test_resource.replication_task_arn
}

output "target_endpoint_arn" {
  value = aws_dms_endpoint.named_test_resource_target.endpoint_arn
}

output "replication_instance_arn" {
  value = aws_dms_replication_instance.named_test_resource.replication_instance_arn
}

output "status" {
  value = aws_dms_replication_task.named_test_resource.status
}


output "resource_name" {
  value = var.resource_name
}
