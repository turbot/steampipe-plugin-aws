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

# Create AWS > ECS > Cluster
resource "aws_ecs_cluster" "named_test_resource" {
  name = var.resource_name
}

# Create AWS > ECS > Task Definition
resource "aws_ecs_task_definition" "named_test_resource" {
  family                = var.resource_name
  container_definitions = <<TASK_DEFINITION
  [
    {
        "cpu": 10,
        "essential": true,
        "image": "jenkins",
        "memory": 128,
        "name": "jenkins"
    }
  ]
  TASK_DEFINITION
}

# Create AWS > ECS > Service
resource "aws_ecs_service" "named_test_resource" {
  name            = var.resource_name
  cluster         = aws_ecs_cluster.named_test_resource.id
  task_definition = aws_ecs_task_definition.named_test_resource.arn
  desired_count   = 3

  ordered_placement_strategy {
    type  = "binpack"
    field = "cpu"
  }

  placement_constraints {
    type       = "memberOf"
    expression = "attribute:ecs.availability-zone in [us-west-2a, us-west-2b]"
  }

  tags = {
    name = var.resource_name
  }
}

output "resource_name" {
  value = var.resource_name
}

output "resource_aka" {
  value = aws_ecs_service.named_test_resource.id
}

output "cluster_arn" {
  value = aws_ecs_service.named_test_resource.cluster
}

output "task_definition_arn" {
  value = aws_ecs_task_definition.named_test_resource.arn
}
