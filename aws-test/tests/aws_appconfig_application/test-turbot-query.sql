select
  title,
  tags,
  akas
from
  aws.aws_appconfig_application
where id = '{{ output.id.value }}';