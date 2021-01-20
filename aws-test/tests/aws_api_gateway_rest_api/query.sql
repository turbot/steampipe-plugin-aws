select name
from aws.aws_api_gateway_rest_api
where name = '{{resourceName}}'
order by name