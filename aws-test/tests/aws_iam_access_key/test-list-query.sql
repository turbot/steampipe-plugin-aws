select user_name, access_key_id, status
from aws.aws_iam_access_key
where access_key_id = '{{ output.resource_id.value }}'
