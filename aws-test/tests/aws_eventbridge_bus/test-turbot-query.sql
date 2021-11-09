select akas, name, region, tags, title
from aws.aws_eventbridge_bus
where name = '{{ resourceName }}';
