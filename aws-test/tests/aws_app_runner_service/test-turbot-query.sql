select title, akas, region, account_id
from aws_app_runner_service
where arn = '{{ output.resource_aka.value }}'
