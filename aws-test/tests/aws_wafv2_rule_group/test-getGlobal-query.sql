select name, arn, id, scope, description, capacity, rules, visibility_config, partition, region, account_id
from aws.aws_wafv2_rule_group
where id = '{{ output.resource_id_global.value }}' and name = '{{ output.resource_name_global.value }}' and scope = 'CLOUDFRONT';