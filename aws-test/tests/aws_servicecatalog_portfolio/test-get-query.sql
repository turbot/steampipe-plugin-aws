select display_name, id, arn, provider_name
from aws_servicecatalog_portfolio
where id = '{{ output.resource_id.value }}';
