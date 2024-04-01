select
  application_version_arn,
  akas,
  region,
  tags,
  title
from
  aws_elastic_beanstalk_application_version
where
  version_label = '{{ output.version_label.value }}::scooby';