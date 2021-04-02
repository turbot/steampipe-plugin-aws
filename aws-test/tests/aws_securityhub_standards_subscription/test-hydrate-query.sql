select standards_status, standards_subscription_arn, standards_input
from aws.aws_securityhub_standards_subscription
where name = "CIS AWS Foundations Benchmark v1.2.0";
