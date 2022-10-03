select
  arn,
  region,
  tags
from
  aws_iam_saml_provider
where akas::text = '["{{ output.resource_aka.value }}"]';
