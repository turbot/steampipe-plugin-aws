select 
  akas,
  tags,
  title
from 
  aws.aws_codeartifact_repository
where
  name = '{{ resourceName }}';