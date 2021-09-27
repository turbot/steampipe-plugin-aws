select title, akas, partition, region, account_id
from aws.aws_identitystore_group
where id = '{{ output.resource_id.value }}'
and identity_store_id = '{{ output.identity_store_id.value }}';
