select title, region
from aws.aws_iam_ssh_public_key
where user_name = '{{ resourceName }}';
