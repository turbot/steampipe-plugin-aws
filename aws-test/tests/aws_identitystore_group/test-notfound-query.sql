select name
from aws.aws_identitystore_group
where id = '0000000000-00000000-0000-0000-0000-000000000000'
and identity_store_id = '{{ output.identity_store_id.value }}';
