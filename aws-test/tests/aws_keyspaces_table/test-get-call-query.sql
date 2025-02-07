select title, akas, table_name
from aws.aws_keyspaces_table
where keyspace_name = '{{resourceName}}' and table_name = '{{resourceName}}';
