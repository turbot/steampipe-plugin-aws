select
  akas,
  tags,
  title
from
  aws.aws_directconnect_connection
where
  connection_id = '{{ output.connection_id.value }}';
