select title, akas, tags
from aws.aws_wafv2_rule_group
where id = '{{ output.resource_id_regional.value }}';