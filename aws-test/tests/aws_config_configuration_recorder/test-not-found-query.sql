select name, role_arn, status, status_recording, title, akas
from aws.aws_config_configuration_recorder
where name = 'dummy-{{ resourceName }}';