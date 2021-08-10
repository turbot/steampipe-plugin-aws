select name, directory_id, arn
from aws.aws_directory_service_directory
where directory_id = '{{ output.resource_id.value }}';
