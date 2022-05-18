select
  name,
  arn,
  group_by_attribute
from
  aws_securityhub_insight
where
  arn='{{ output.insight_arn.value }}';