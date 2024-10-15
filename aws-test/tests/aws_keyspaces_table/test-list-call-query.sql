select table_name, keyspace_name, title, akas
from aws.aws_keyspaces_table
where keyspace_name = '{{ resourceName }}';
