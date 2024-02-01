select akas, tags, title
from aws.aws_iot_thing_type
where arn = '{{ output.resource_aka.value }}';