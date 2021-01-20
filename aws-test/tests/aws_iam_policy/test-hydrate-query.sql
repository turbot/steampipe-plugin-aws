select name, policy, policy_std, akas
from aws.aws_iam_policy
where arn = '{{ output.resource_aka.value }}'
