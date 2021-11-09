select arn, name, tags_src
from aws.aws_eventbridge_bus
where name = '{{ resourceName }}';
