select name, arn
from aws_appsync_graphql_api
where arn = '{{ output.resource_aka.value }}';
