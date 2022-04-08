select title, akas
from aws_wafregional_rule
where rule_id = '{{ output.rule_id.value }}'; 