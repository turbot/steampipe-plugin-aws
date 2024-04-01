select akas, title
from aws.aws_iot_thing
where thing_name = '{{ resourceName }}';