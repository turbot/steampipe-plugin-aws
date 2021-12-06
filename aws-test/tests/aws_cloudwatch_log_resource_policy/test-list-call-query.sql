select account_id, policy_name, partition, region
from aws.aws_cloudwatch_log_resource_policy
where policy_name = '{{ resourceName }}';
