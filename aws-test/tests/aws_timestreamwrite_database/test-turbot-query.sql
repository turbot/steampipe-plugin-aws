select akas, title
from aws.aws_timestreamwrite_database
where akas::text = '["{{ output.resource_aka.value }}"]';
