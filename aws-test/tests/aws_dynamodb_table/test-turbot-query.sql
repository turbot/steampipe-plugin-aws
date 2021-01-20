select title, akas, region, partition, account_id
from aws.aws_dynamodb_table
where name = '{{ resourceName }}'