select name, arn, role_arn, status, status_recording, title, akas, region
from aws.aws_config_configuration_recorder
where name = '{{ resourceName }}' and  region = '{{ output.region_name.value }}';