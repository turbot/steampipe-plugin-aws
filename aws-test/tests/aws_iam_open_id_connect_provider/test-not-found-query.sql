select
  akas,
  region,
  tags
from
  aws_iam_open_id_connect_provider
where arn = '{{ output.resource_aka.value }}dummy';
