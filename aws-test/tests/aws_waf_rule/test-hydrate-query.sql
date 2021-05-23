select metric_name, name, partition, tags
from aws.aws_waf_rule
where rule_id = '{{ output.rule_id.value }}';