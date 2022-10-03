select
  arn,
  id,
  name,
  description,
  tags
from
  aws_appconfig_application
where
  id = '{{ output.id.value }}';