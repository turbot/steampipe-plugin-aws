select name, conformance_pack_id, conformance_pack_input_parameters, akas
from aws.aws_config_conformance_pack
where name = 'dummy-{{ resourceName }}';