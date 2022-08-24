select
  arn,
  id,
  name,
  description
from
  aws.aws_appconfig_application
where arn = '{{ output.resource_aka.value }}';