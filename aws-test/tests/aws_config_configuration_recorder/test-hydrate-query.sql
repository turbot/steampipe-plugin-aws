select name, role_arn, status, status_recording
from aws.aws_config_configuration_recorder
where name = '{{ resourceName }}'
