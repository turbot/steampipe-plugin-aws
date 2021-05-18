select name, arn, id, scope, description, regular_expressions, partition, region, account_id
from aws.aws_wafv2_regex_pattern_set
where id = '{{ output.resource_id_regional.value }}' and name = '{{ output.resource_name_regional.value }}' and scope = 'REGIONAL';