select user_name, ssh_public_key_id
from aws.aws_iam_ssh_public_key
where user_name = '{{ resourceName }}';
