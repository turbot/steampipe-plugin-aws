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

resource "aws_iam_instance_profile" "instance_profile" {
  name = var.resource_name
  role = aws_iam_role.role.name
}

resource "aws_iam_role" "role" {
  name = var.resource_name
  path = "/"

  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Principal": {
               "Service": "ec2.amazonaws.com"
            },
            "Effect": "Allow",
            "Sid": ""
        }
    ]
}
EOF
}

resource "aws_vpc" "my_vpc" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_internet_gateway" "testinternetgateway" {
  vpc_id = aws_vpc.my_vpc.id
}

resource "aws_route_table" "route" {
  vpc_id = "${aws_vpc.my_vpc.id}"
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.testinternetgateway.id}"
  }
}

resource "aws_subnet" "my_subnet1" {
  vpc_id            = "${aws_vpc.my_vpc.id}"
  availability_zone = "${var.aws_region}b"
  map_public_ip_on_launch = "true"
  cidr_block        = "10.0.2.0/24"
}

resource "aws_route_table_association" "routeassociation1" {
  subnet_id      = aws_subnet.my_subnet1.id
  route_table_id = aws_route_table.route.id
}

resource "aws_security_group" "ssh-allowed" {
    vpc_id = aws_vpc.my_vpc.id
    egress {
        from_port = 0
        to_port = 0
        protocol = -1
        cidr_blocks = ["0.0.0.0/0"]
    }
    ingress {
        from_port = 22
        to_port = 22
        protocol = "tcp"
        // This means, all ip address are allowed to ssh !
        // Do not do it in the production.
        // Put your office or home address in it!
        cidr_blocks = ["0.0.0.0/0"]
    }
    //If you do not add this rule, you can not reach the NGIX
    ingress {
        from_port = 80
        to_port = 80
        protocol = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }
}

resource "aws_elastic_beanstalk_application" "application_test" {
  name = var.resource_name
  appversion_lifecycle {
    service_role = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/aws-service-role/elasticbeanstalk.amazonaws.com/AWSServiceRoleForElasticBeanstalk"
  }
}

resource "aws_elastic_beanstalk_environment" "named_test_resource" {
  name                   = var.resource_name
  application            = aws_elastic_beanstalk_application.application_test.name
  solution_stack_name    = "64bit Amazon Linux 2 v3.6.4 running Go 1"
  wait_for_ready_timeout = "45m"
  setting {
    namespace = "aws:autoscaling:launchconfiguration"
    name      = "IamInstanceProfile"
    value     = aws_iam_instance_profile.instance_profile.name
  }
  setting {
    namespace = "aws:ec2:vpc"
    name      = "VPCId"
    value     = "${aws_vpc.my_vpc.id}"
  }
  setting {
    namespace = "aws:ec2:vpc"
    name      = "ELBSubnets"
    value     = "${aws_subnet.my_subnet1.id}"
  }
  setting {
    namespace = "aws:ec2:vpc"
    name      = "Subnets"
    value     = "${aws_subnet.my_subnet1.id}"
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
  value = aws_elastic_beanstalk_environment.named_test_resource.arn
}

output "resource_id" {
  value = aws_elastic_beanstalk_environment.named_test_resource.id
}

output "resource_name" {
  value = var.resource_name
}
