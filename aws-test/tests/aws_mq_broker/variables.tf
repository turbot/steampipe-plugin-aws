
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

resource "aws_mq_configuration" "named_test_resource" {
  description    = "Example Configuration"
  name           = var.resource_name
  engine_type    = "ActiveMQ"
  engine_version = "5.17.6"

  data = <<DATA
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<broker xmlns="http://activemq.apache.org/schema/core">
  <plugins>
    <forcePersistencyModeBrokerPlugin persistenceFlag="true"/>
    <statisticsBrokerPlugin/>
    <timeStampingBrokerPlugin ttlCeiling="86400000" zeroExpirationOverride="86400000"/>
  </plugins>
</broker>
DATA
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
  tags = {
    name = var.resource_name
  }
}

resource "aws_security_group" "named_test_resource" {
  vpc_id      = aws_vpc.main.id
  name        = var.resource_name
  description = "Test Security Group."
  tags = {
    name = var.resource_name
  }
}

resource "aws_subnet" "named_test_resource" {
  vpc_id     = aws_vpc.main.id
  cidr_block = "10.0.1.0/24"
  availability_zone = "us-east-1a" 
  tags = {
    Name = var.resource_name
  }
}

resource "aws_mq_broker" "named_test_resource" {
  depends_on  = [aws_subnet.named_test_resource, aws_security_group.named_test_resource, aws_vpc.main, aws_mq_configuration.named_test_resource]
  broker_name = var.resource_name

  configuration {
    id       = aws_mq_configuration.named_test_resource.id
    revision = aws_mq_configuration.named_test_resource.latest_revision
  }

  engine_type        = "ActiveMQ"
  engine_version     = "5.17.6"
  host_instance_type = "mq.m5.large"
  security_groups    = [aws_security_group.named_test_resource.id]
  subnet_ids = [aws_subnet.named_test_resource.id]

  user {
    username = "ExampleUser"
    password = "MindTheGapGood"
  }

  tags = {
    Name = var.resource_name
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

output "resource_arn" {
  depends_on  = [aws_mq_broker.named_test_resource]
  value       = aws_mq_broker.named_test_resource.arn
}

output "resource_id" {
  depends_on  = [aws_mq_broker.named_test_resource]
  value       = aws_mq_broker.named_test_resource.id
}

output "engine_type" {
  value = "ActiveMQ"
}

output "engine_version" {
  value = "5.17.6"
}

output "resource_name" {
  value = var.resource_name
}
