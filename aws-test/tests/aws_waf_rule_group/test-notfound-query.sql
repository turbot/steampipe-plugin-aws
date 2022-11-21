select
  name,
  arn,
  rule_group_id,
  metric_name,
  tags
from
  aws_waf_rule_group
where 
  rule_group_id = '{{ output.rule_group_id.value }}-dummy';