select
  title,
  tags,
  akas
from
  aws_appconfig_application
where
  id = '{{ output.id.value }}';