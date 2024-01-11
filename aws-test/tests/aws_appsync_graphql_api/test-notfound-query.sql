select name, arn
from aws_appsync_graphql_api
where name = 'dummy-{{ resourceName }}';
