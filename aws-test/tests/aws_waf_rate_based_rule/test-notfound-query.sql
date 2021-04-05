select metric_name, region, account_id
from aws.aws_waf_rate_based_rule
where rule_id = '{{ output.resource_id.value }}.1ad';
