select name
from aws.aws_api_gateway_api_key
where name = '{{resourceName}}'
order by name
