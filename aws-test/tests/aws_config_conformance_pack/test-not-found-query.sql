select name, conformance_pack_id, input_parameters, created_by_svc, delivery_bucket, delivery_bucket_prefix, last_update, title, akas
from aws.aws_config_conformance_pack
where name = 'dummy-{{ resourceName }}';