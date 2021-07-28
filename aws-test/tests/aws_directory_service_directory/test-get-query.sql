select name, directory_id, arn
from aws.aws_directory_service_directory
where name = '{{ output.resource_name.value }}.com';
