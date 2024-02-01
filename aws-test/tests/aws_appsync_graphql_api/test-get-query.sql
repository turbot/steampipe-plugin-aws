select name, arn, api_id, api_type
from aws_appsync_graphql_api
where api_id = '{{ output.resource_id.value }}'
