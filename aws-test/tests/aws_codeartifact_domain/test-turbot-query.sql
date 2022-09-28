select 
  akas,
  tags,
  title
from 
  aws.aws_codeartifact_domain
where
  name = '{{ resourceName }}';