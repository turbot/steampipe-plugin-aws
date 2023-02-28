
variable "resource_name" {
  type        = string
  default     = "turbot-test-20221123-create-update"
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
    Name = var.resource_name
  }
}

resource "aws_security_group" "security_group" {
  name   = var.resource_name
  vpc_id = aws_vpc.my_vpc.id
  tags = {
    Name = var.resource_name
  }
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.my_vpc.id
  tags = {
    Name = var.resource_name
  }
}

resource "aws_subnet" "my_subnet1" {
  cidr_block        = "10.1.1.0/24"
  availability_zone = "${var.aws_region}a"
  vpc_id            = aws_vpc.my_vpc.id
  tags = {
    Name = var.resource_name
  }
  depends_on = [
    aws_internet_gateway.igw
  ]
}

resource "aws_subnet" "my_subnet2" {
  cidr_block        = "10.1.2.0/24"
  availability_zone = "${var.aws_region}b"
  vpc_id            = aws_vpc.my_vpc.id
  tags = {
    Name = var.resource_name
  }
  depends_on = [
    aws_internet_gateway.igw
  ]
}

resource "aws_iam_role" "role" {
  name = var.resource_name
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole",
        Effect = "Allow",
        Principal = {
          Service = "rds.amazonaws.com"
        }
      }
    ]
  })
  inline_policy {
    name = "${var.resource_name}_inline_policy"
    policy = jsonencode({
      Version = "2012-10-17"
      Statement = [
        {
          "Sid" : "GetSecretValue",
          "Action" : [
            "secretsmanager:GetSecretValue"
          ],
          "Effect" : "Allow",
          "Resource" : ["*"]
        },
        {
          "Sid" : "DecryptSecretValue",
          "Action" : [
            "kms:Decrypt"
          ],
          "Effect" : "Allow",
          "Resource" : ["*"],
          "Condition" : {
            "StringEquals" : {
              "kms:ViaService" : "secretsmanager.${var.aws_region}.amazonaws.com"
            }
          }
        }
      ]
    })
  }
}

resource "random_integer" "rand_int" {
  min = 1
  max = 50000
  keepers = {
    listener_arn = aws_vpc.my_vpc.id
  }
}

resource "aws_secretsmanager_secret" "secret" {
  name = "${var.resource_name}-${random_integer.rand_int.result}"
}

resource "aws_db_proxy" "named_test_resource" {
  name                   = var.resource_name
  debug_logging          = false
  engine_family          = "MYSQL"
  idle_client_timeout    = 1800
  require_tls            = true
  role_arn               = aws_iam_role.role.arn
  vpc_security_group_ids = [aws_security_group.security_group.id]
  vpc_subnet_ids         = [aws_subnet.my_subnet1.id, aws_subnet.my_subnet2.id]

  auth {
    auth_scheme = "SECRETS"
    description = "example"
    iam_auth    = "DISABLED"
    secret_arn  = aws_secretsmanager_secret.secret.arn
  }

  tags = {
    Name = var.resource_name
  }
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

output "resource_id" {
  value = aws_db_proxy.named_test_resource.id
}

output "resource_aka" {
  value = aws_db_proxy.named_test_resource.arn
}
