select
  document_version,
  document_format,
  document_type,
  name,
  partition, region, tags, title
from
  aws.aws_ssm_document
where
  akas::text = '["{{ output.resource_aka.value }}"]';

