select akas, title
from aws.aws_iot_thing_group
where group_name = '{{ resourceName }}';