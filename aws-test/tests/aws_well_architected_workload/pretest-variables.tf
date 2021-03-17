
variable "resource_name" {
  type    = string
  default = ""
  description = "Name of the resource used throughout the test."
}
variable "resource_name_1" {
  type    = string
  default = ""
  description = "Name of the resource used throughout the test."
}
variable "resource_name_2" {
  type    = string
  default = ""
  description = "Name of the resource used throughout the test."
}

variable "turbot_profile" {
  type    = string
  default = "osbornNew"
  description = "Turbot credentials profile to use for the test run."
}

provider "turbot" {
  profile = var.turbot_profile
}


variable "aws_profile" {
  type    = string
  default = "osbornNew"
  description = "AWS credentials profile used for the test. Default is to use the default profile."
}

variable "aws_region" {
  type    = string
  default = "us-east-2"
  description = "AWS region used for the test. Does not work with default region in config, so must be defined here."
}

variable "aws_region_alternate" {
  type    = string
  default = "us-east-1"
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


resource "null_resource" "named_test_resource" {
  provisioner "local-exec" {
    command = <<EOT
      aws wellarchitected create-workload --workload-name ${var.resource_name} --description 'testinh-cli'  --environment 'PREPRODUCTION' --review-owner 'khushboo' --lenses 'wellarchitected' --aws-regions 'us-east-1' > ${data.template_file.log_name.rendered};
    EOT
  }
}

data "template_file" "log_name" {
  template = "${path.module}/output.log"
}

data "local_file" "named_test_resource" {
  filename = "${data.template_file.log_name.rendered}"
  depends_on = ["null_resource.named_test_resource"]
}

output "resource_id" {
  value = "${data.local_file.named_test_resource.content}"
}


# Get the Turbot resource that represents the test resource. This shadow will be used
# for all policy settings etc throughout the test. Waiting for the shadow also ensures
# proper ordering of the setup of policies etc after the resource exists in Turbot.

resource "turbot_shadow_resource" "shadow_resource" {
  resource = "arn:${data.aws_partition.current.partition}::${var.aws_region}:${data.aws_caller_identity.current.account_id}"
}


output "resource_name" {
  value = var.resource_name
}


