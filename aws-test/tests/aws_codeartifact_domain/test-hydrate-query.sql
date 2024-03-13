select 
  arn,
  name, 
  owner,
  repository_count::text,
  asset_size_bytes::text, 
  tags
from 
  aws.aws_codeartifact_domain
where
  arn = '{{ output.resource_aka.value }}';
