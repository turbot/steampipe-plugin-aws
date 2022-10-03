select
  name,
  arn,
  rule_group_id,
  metric_name,
  tags
from
  aws_waf_rule_group
where 
  akas::text = '["{{ output.resource_aka.value }}"]';