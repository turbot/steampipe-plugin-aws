select title, rule_arn, tags, akas 
from aws.aws_config_rule
where name = '{{ resourceName }}';