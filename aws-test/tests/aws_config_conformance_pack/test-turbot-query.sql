select name, title, account_id, region, akas
from aws.aws_config_conformance_pack
where name = '{{ resourceName }}';