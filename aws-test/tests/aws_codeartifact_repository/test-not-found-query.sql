select
  arn,
  administrator_account,
  description,
  domain_owner,
  region,
  tags
from
  aws_codeartifact_repository
where
  name = 'dummy-{{ resourceName }}';
