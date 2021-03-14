select name, conformance_pack_id, input_parameters, created_by_svc, delivery_bucket, delivery_bucket_prefix, last_update, title
from aws.aws_config_conformance_pack
where akas::text = '["{{ output.resource_aka.value }}"]';
