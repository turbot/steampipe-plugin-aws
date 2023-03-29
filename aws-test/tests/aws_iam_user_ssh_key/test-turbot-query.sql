select title, region
from aws.aws_iam_user_ssh_key
where user_name = '{{ resourceName }}';
