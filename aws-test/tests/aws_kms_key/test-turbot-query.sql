select akas, tags, title
from aws.aws_kms_key
where id='{{ output.resource_id.value }}'
