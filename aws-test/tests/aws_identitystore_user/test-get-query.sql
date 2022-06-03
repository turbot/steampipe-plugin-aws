select name
from aws.aws_identitystore_user
where id = '{{ output.resource_id.value }}'
and identity_store_id = '{{ output.identity_store_id.value }}';
