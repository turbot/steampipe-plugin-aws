select volume_id, auto_enable_io, product_codes
from aws.aws_ebs_volume
where volume_id = '{{ output.resource_id.value }}'
