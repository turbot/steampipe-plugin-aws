select server_id, arn, user_name
from aws.aws_transfer_user
where server_id = '{{output.resource_server_id.value}}'
