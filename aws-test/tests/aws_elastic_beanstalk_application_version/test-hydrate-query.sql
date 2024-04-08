select
  application_version_arn,
  version_label,
  tags_src
from
  aws_elastic_beanstalk_application_version
where
  application_version_arn = '{{ output.application_version_arn.value }}';