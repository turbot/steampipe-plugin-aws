select *
from aws.aws_iot_thing_group
where group_name = 'dummy-{{ resourceName }}';