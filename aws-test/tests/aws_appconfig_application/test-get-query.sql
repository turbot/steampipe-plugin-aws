select
  arn,
  id,
  name,
  description
from
  aws.aws_appconfig_application
where id = '{{ output.id.value }}';