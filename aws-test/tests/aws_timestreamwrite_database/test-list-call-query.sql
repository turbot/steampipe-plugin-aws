select arn, database_name
from aws.aws_timestreamwrite_database
where akas::text = '["{{ output.resource_aka.value }}"]';
