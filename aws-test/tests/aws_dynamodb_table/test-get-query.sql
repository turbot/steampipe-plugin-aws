select name, arn
from aws.aws_dynamodb_table
where name = '{{ resourceName }}'
