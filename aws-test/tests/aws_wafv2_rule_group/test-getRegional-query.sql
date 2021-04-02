select name, arn, id, scope, description, capacity, rules, visibility_config, partition, region, account_id
from aws.aws_wafv2_rule_group
where id = '{{ output.resource_id_regional.value }}';