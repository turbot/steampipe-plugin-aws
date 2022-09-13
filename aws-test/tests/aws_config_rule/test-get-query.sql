select name, rule_id, arn, rule_state, tags_src, description, title, akas 
from aws_config_rule
where name = '{{ resourceName }}';