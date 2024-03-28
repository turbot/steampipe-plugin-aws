variable "resource_name" {
  type        = string
  default     = "turbot-test-28032024"
  description = "Name of the resource used throughout the test."
}

variable "application_name" {
  type        = string
  default     = "tf-test-app"
  description = "Name of the Elastic Beanstalk application used for testing."
}

variable "version_label" {
  type        = string
  default     = "tf-test-app-version-28032024"
  description = "Version label for the application version."
}

variable "aws_profile" {
  type        = string
  default     = "default"
  description = "AWS credentials profile used for the test. Default is to use the default profile."
}

variable "aws_region" {
  type        = string
  default     = "us-east-1"
  description = "AWS region used for the test."
}

provider "aws" {
  profile = var.aws_profile
  region  = var.aws_region
}

resource "aws_s3_bucket" "application_version_bucket" {
  bucket = "tf-test-application-version-bucket-${var.application_name}"
}

resource "aws_s3_object" "application_version_object" {
  bucket = aws_s3_bucket.application_version_bucket.id
  key    = "beanstalk/go.zip"
  source = "/Users/temp/Downloads/go.zip" # Path to the file to be uploaded to S3 update the path as per your local file path
}

resource "aws_elastic_beanstalk_application" "test_application" {
  name        = var.application_name
  description = "An application to test AWS Elastic Beanstalk Application Version."
}

resource "aws_elastic_beanstalk_application_version" "test_application_version" {
  name             = var.version_label
  application      = aws_elastic_beanstalk_application.test_application.name
  description      = "A version of the test application."
  bucket           = aws_s3_bucket.application_version_bucket.id
  key              = aws_s3_object.application_version_object.key
  tags = {
    name = var.version_label
  }
}

output "aws_region" {
  value = var.aws_region
}

output "application_name" {
  value = aws_elastic_beanstalk_application.test_application.name
}

output "version_label" {
  value = aws_elastic_beanstalk_application_version.test_application_version.name
}

output "application_version_arn" {
  value = aws_elastic_beanstalk_application_version.test_application_version.arn
}
