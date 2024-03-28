select
  application_version_arn,
  application_name,
  version_label,
  tags
from
  aws_elastic_beanstalk_application_version
where
  application_name = '{{ output.application_name.value }}' and version_label = '{{ output.version_label.value }}';