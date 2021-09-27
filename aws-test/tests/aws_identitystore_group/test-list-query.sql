select id
from aws.aws_identitystore_group
where name = '{{ output.resource_name.value }}'
and identity_store_id = '{{ output.identity_store_id.value }}';
