select metric_name, akas, title, tags
from aws.aws_waf_rate_based_rule
where rule_id = '{{ output.resource_id.value }}';

