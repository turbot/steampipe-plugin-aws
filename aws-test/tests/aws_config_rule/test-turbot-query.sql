select title, arn, tags, akas 
from aws_config_rule
where name = '{{ resourceName }}';