select title, partition, region, account_id
from aws.aws_cloudwatch_log_resource_policy
where policy_name = '{{ resourceName }}';
