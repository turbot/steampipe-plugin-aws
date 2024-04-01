select *
from aws.aws_iot_thing
where thing_name = 'dummy-{{ resourceName }}';