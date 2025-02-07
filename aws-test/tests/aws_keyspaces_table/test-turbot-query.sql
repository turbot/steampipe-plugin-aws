select akas, title
from aws.aws_keyspaces_table
where keyspace_name = '{{ resourceName }}';
