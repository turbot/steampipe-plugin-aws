select id, arn, name
from aws_pinpoint_app
where id = 'dummy-{{ resourceName }}';
