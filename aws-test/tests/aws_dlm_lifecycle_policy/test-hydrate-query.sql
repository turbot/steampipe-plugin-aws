select policy_id, arn
from aws.aws_dlm_lifecycle_policy
where policy_id = '{{ output.resource_id.value }}'
