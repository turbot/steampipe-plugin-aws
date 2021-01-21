
select id
from aws.aws_kms_key
where id='{{ output.resource_id.value }}'
