select
  name,
  arn,
  group_by_attribute
from
  aws_securityhub_insight
where 
  name = '{{ resourceName }}';