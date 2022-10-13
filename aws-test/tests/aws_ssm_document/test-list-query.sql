select
  document_version,
  document_format,
  document_type,
  name,
  partition, region, tags, title
from
  aws.aws_ssm_document
where
  owner = 'Self'
  and akas::text = '["{{ output.resource_aka.value }}"]';

