select standards_arn, name, enabled_by_default, standards_status, standards_subscription_arn
from aws.aws_securityhub_standards_subscription
where standards_arn = "{{ output.cis_arn.value }};
