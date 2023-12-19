select
  arn,
  region,
  tags
from
  aws_iam_open_id_connect_provider
where akas::text = '["{{ output.resource_aka.value }}"]';
