select name, table_arn
from aws.aws_dynamodb_table
where name = '{{ resourceName }}'
