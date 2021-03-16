select name, title, account_id, region, akas
from aws.aws_config_configuration_recorder
where name = '{{ resourceName }}';