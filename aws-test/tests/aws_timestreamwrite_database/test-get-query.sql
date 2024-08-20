select arn, database_name
from aws.aws_timestreamwrite_database
where database_name = '{{ resourceName }}';
