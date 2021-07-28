select name, akas, tags, title
from aws.aws_directory_service_directory
where name = '{{ output.resource_name.value }}';
