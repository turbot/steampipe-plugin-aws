select title, partition, region, account_id
from aws.aws_ssoadmin_managed_policy_attachment
where permission_set_arn = '{{ output.resource_arn.value }}';
