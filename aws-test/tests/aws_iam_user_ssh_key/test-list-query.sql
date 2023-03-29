select user_name, ssh_public_key_id
from aws.aws_iam_user_ssh_key
where user_name = '{{ resourceName }}';
