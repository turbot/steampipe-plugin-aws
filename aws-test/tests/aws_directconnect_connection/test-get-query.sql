select connection_id, connection_name, connection_state, bandwidth, location, arn, tags_src
from aws.aws_directconnect_connection
where connection_id = '{{ output.connection_id.value }}'
