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

resource "aws_vpc" "main" {
  cidr_block       = "10.0.0.0/16"
  instance_tenancy = "default"
  tags = {
    Name = "demo-for-int-test"
  }
}

# Creating subnet
resource "aws_subnet" "demosubnet" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = "us-east-1a"
  map_public_ip_on_launch = true
  tags = {
    Name = "Public Subnet"
  }
}

# Creating Internet Gateway
resource "aws_internet_gateway" "demogateway" {
  vpc_id = aws_vpc.main.id
}

# Creating Route Table for Public Subnet
resource "aws_route_table" "rt" {
  vpc_id = aws_vpc.main.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.demogateway.id
  }
  tags = {
    Name = "Public Subnet Route Table"
  }
}
resource "aws_route_table_association" "rt_associate_public" {
  subnet_id      = aws_subnet.demosubnet.id
  route_table_id = aws_route_table.rt.id
}

# Creating Security Group
resource "aws_security_group" "demosg" {
  vpc_id      = "${aws_vpc.main.id}"
  # Outbound Rules
  # Internet access to anywhere
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

data "aws_ami" "linux" {
  most_recent = true
  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-2.0.20220606.1-x86_64-gp2"]
  }
  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
  owners = ["137112412989"]
}

# Instance Profile
resource "aws_iam_instance_profile" "test_profile" {
  name = "test_profile"
  role = aws_iam_role.test_role.name
}

# IAM Role for Instance Profile
resource "aws_iam_role" "test_role" {
  name = "test_role"
  path = "/"
  # Terraform expression result to valid JSON syntax.
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      },
    ]
  })
}
# Attach AmazonSSMManagedInstanceCore Policy to the role for SSM
resource "aws_iam_role_policy_attachment" "test-attach" {
  role       = aws_iam_role.test_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

# Creating EC2 instance in Public Subnet
resource "aws_instance" "named_test_resource" {
  ami           = data.aws_ami.linux.id
  instance_type = "t2.micro"
  vpc_security_group_ids = [ aws_security_group.demosg.id ]
  subnet_id = aws_subnet.demosubnet.id
  associate_public_ip_address = true
  iam_instance_profile = aws_iam_instance_profile.test_profile.name
  tags = {
    Name = var.resource_name
  }
}

resource "null_resource" "delay" {
  depends_on = [
    aws_instance.named_test_resource
  ]
  provisioner "local-exec" {
    command = "sleep 60"
  }
}

resource "null_resource" "create_association" {
  depends_on = [null_resource.delay]
  provisioner "local-exec" {
    command = <<EOF
      aws ssm create-association --name "AWS-UpdateSSMAgent" --targets "Key=instanceids,Values=${aws_instance.named_test_resource.id}" --region ${data.aws_region.primary.name} --profile ${var.aws_profile}
    EOF
  }
}

resource "null_resource" "delay_create_association" {
  depends_on = [
    aws_instance.named_test_resource
  ]
  provisioner "local-exec" {
    command = "sleep 300"
  }
}

locals {
  path = "${path.cwd}/compliance.json"
}

resource "null_resource" "list_compliance" {
  depends_on = [null_resource.delay_create_association]
  provisioner "local-exec" {
    command = "aws ssm list-compliance-items --resource-ids ${aws_instance.named_test_resource.id} --resource-types ManagedInstance --output json --profile ${var.aws_profile} --region ${data.aws_region.primary.name} > ${local.path}"
  }
}

data "local_file" "input" {
  depends_on = [null_resource.list_compliance]
  filename   = local.path
}

output "arn" {
  depends_on = [null_resource.list_compliance]
  value      = "arn:${data.aws_partition.current.partition}:ssm:${data.aws_region.primary.name}:${data.aws_caller_identity.current.account_id}:managed-instance/${aws_instance.named_test_resource.id}/compliance-item/${jsondecode(data.local_file.input.content).ComplianceItems[0].Id}:${jsondecode(data.local_file.input.content).ComplianceItems[0].ComplianceType}"
}

output "resource_id" {
  value = aws_instance.named_test_resource.id
}

output "compliance_id" {
  depends_on = [null_resource.list_compliance]
  value      = jsondecode(data.local_file.input.content).ComplianceItems[0].Id
}

output "status" {
  depends_on = [null_resource.list_compliance]
  value      = jsondecode(data.local_file.input.content).ComplianceItems[0].Status
}
