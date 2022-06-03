select policy_id, arn
from aws.aws_dlm_lifecycle_policy
where arn = '{{ output.resource_aka.value }}'
