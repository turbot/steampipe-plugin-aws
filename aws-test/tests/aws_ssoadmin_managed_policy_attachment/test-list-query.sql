select permission_set_arn, managed_policy_arn
from aws.aws_ssoadmin_managed_policy_attachment
where permission_set_arn = '{{ output.resource_arn.value }}';
