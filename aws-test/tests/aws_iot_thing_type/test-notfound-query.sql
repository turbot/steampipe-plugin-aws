select *
from aws.aws_iot_thing_type
where thing_type_name = 'dummy{{ resourceName }}';