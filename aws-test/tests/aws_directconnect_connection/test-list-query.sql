select connection_id, connection_name, arn
from aws.aws_directconnect_connection
where connection_id = '{{ output.connection_id.value }}'
