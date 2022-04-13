select name, rule_id
from aws_wafregional_rule
where rule_id = '{{ output.rule_id.value }}';