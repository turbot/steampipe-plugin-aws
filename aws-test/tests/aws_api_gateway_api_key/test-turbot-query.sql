select akas, title, tags
from aws.aws_api_gateway_api_key
where id = '{{ output.resource_id.value }}';
