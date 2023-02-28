select
  akas,
  tags,
  title
from
  aws_codeartifact_repository
where
  name = '{{ resourceName }}';