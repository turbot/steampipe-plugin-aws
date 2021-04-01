select name, conformance_pack_id, input_parameters, akas
from aws.aws_config_conformance_pack
where name = 'dummy-{{ resourceName }}';