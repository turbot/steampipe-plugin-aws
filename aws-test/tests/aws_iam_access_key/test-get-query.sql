select user_name, access_key_id, status
from aws_iam_access_key
where user_name = '{{ output.user_name.value }}'
