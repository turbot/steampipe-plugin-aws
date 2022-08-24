select
  arn,
  id,
  name,
  description,
  tags
from
  aws.aws_appconfig_application
where arn = '{{ output.resource_aka.value }}';