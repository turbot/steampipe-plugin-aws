select
  name,
  insight_arn,
  group_by_attribute
from
  aws_securityhub_insight
where
  insight_arn='{{ output.insight_arn.value }}';