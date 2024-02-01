select thing_type_name, arn, account_id, partition, region
from aws.aws_iot_thing_type
where thing_type_name = '{{ resourceName }}';