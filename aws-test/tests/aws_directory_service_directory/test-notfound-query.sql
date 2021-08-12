select title, akas, region, account_id
from aws.aws_directory_service_directory
where name = 'dummy-{{ output.resource_name.value }}';