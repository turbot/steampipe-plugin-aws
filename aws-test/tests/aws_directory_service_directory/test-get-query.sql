select name, directory_id, arn, description, access_url
from aws.aws_directory_service_directory
where cluster_name = '{{ output.resource_name.value }}';
