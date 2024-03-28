select
  application_version_arn,
  version_label,
  region,
  tags,
  title
from
  aws_elastic_beanstalk_application_version
where
  akas:: text = '["{{ output.application_version_arn.value }}"]';