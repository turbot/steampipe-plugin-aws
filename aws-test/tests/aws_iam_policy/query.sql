select name, arn, policy, policy_std
from aws.aws_iam_policy
where arn = '{{ output.resource_aka.value }}'
