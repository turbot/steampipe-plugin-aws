select name, conformance_pack_input_parameters, delivery_s3_bucket, delivery_s3_key_prefix, title, akas
from aws.aws_config_conformance_pack
where name = '{{ resourceName }}';