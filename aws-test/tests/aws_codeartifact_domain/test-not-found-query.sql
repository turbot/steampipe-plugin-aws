select 
  arn,
  name, 
  owner,
  repository_count,
  asset_size_bytes, 
  tags_src
from 
  aws.aws_codeartifact_domain
where
  arn = 'dummy-{{ output.resource_aka.value }}';
