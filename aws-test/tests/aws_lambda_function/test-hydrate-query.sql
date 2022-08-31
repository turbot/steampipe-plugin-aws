select tags, policy, policy_std, reserved_concurrent_executions
from aws.aws_lambda_function
where name = '{{ resourceName }}';
