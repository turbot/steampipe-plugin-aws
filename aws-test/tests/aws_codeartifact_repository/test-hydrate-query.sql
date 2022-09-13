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
  arn = '{{ output.resource_aka.value }}';
