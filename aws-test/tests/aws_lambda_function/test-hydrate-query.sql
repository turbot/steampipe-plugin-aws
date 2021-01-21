select tags, policy, policy_std
from aws.aws_lambda_function
where name = '{{ resourceName }}'
