select akas
from aws.aws_transfer_server
where server_id = '{{output.resource_id.value}}'
