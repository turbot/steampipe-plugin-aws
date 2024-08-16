select arn, table_name
from aws.aws_timestreamwrite_table
where akas::text = '["{{ output.resource_aka.value }}"]';
