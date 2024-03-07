select thing_name, arn, version::text, attributes
from aws.aws_iot_thing
where title = '{{ resourceName }}';