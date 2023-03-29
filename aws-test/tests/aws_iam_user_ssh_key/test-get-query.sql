select ssh_public_key_id, user_name
from aws.aws_iam_user_ssh_key
where user_name = '{{ resourceName }}' and ssh_public_key_id = '{{ output.ssh_public_key_id.value }}';
