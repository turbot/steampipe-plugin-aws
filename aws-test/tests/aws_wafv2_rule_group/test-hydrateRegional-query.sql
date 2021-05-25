select name, description, tags_src
from aws.aws_wafv2_rule_group
where id = '{{ output.resource_id_regional.value }}';