select akas
from aws.aws_transfer_user
where server_id = '{{output.resource_server_id.value}}'
and user_name = '{{output.resource_user_name.value}}'
