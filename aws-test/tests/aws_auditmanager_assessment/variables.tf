
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

// Create S3 bucket to store assessment reports
resource "aws_s3_bucket" "named_test_resource" {
  bucket        = var.resource_name
  force_destroy = true
}

// Get the assessment frameworks
resource "null_resource" "list_frameworks" {
  provisioner "local-exec" {
    command = "aws auditmanager list-assessment-frameworks --framework-type Standard --output json --profile ${var.aws_profile} --region ${data.aws_region.primary.name} > ${path.cwd}/list-framework.json"
  }
}

data "local_file" "framework" {
  depends_on = [null_resource.list_frameworks]
  filename   = "${path.cwd}/list-framework.json"
}

resource "null_resource" "named_test_resource" {
  depends_on = [
    aws_s3_bucket.named_test_resource,
    null_resource.list_frameworks
  ]
  provisioner "local-exec" {
    command = "aws auditmanager create-assessment --name ${var.resource_name} --description 'Test assessment to validate table outcomes.' --scope 'awsAccounts=[{id=\"${data.aws_caller_identity.current.account_id}\"}],awsServices=[{serviceName=\"ec2\"}]' --roles 'roleArn=\"${data.aws_caller_identity.current.arn}\",roleType=\"PROCESS_OWNER\"' --assessment-reports-destination 'destinationType=\"S3\",destination=\"s3://${var.resource_name}\"' --framework-id '${jsondecode(data.local_file.framework.content).frameworkMetadataList[0].id}' --profile ${var.aws_profile} --region  ${data.aws_region.primary.name} > ${path.cwd}/assessment.json"
  }
}

resource "null_resource" "tag_assessment" {
  depends_on = [null_resource.named_test_resource]
  provisioner "local-exec" {
    command = "aws auditmanager tag-resource --resource-arn '${jsondecode(data.local_file.assessment.content).assessment.arn}' --tags name=${var.resource_name} --profile ${var.aws_profile} --region ${data.aws_region.primary.name}"
  }
}

data "local_file" "assessment" {
  depends_on = [null_resource.named_test_resource]
  filename   = "${path.cwd}/assessment.json"
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "assessment_id" {
  value = jsondecode(data.local_file.assessment.content).assessment.metadata.id
}

output "assessment_arn" {
  value = jsondecode(data.local_file.assessment.content).assessment.arn
}

output "aws_region" {
  value = data.aws_region.primary.name
}

output "aws_partition" {
  value = data.aws_partition.current.partition
}

output "resource_name" {
  value = var.resource_name
}
