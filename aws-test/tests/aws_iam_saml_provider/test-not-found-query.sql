select
  akas,
  region,
  tags
from
  aws_iam_saml_provider
where arn = '{{ output.resource_aka.value }}dummy';
