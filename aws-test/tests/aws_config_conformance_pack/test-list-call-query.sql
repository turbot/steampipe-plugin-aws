select name, conformance_pack_input_parameters, title
from aws.aws_config_conformance_pack
where akas::text = '["{{ output.resource_aka.value }}"]';
