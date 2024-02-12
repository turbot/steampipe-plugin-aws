select thing_name, arn, default_client_id
from aws.aws_iot_thing
where thing_name = '{{ resourceName }}';