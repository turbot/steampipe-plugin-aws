select table_name, arn, database_name
from aws.aws_timestreamwrite_table
where table_name = '{{ resourceName }}' and database_name = '{{ resourceName }}';
