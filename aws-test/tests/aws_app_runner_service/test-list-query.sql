select arn, service_name, service_id
from aws_app_runner_service
where akas::text = '["{{ output.resource_aka.value }}"]'
