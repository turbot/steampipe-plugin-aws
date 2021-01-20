select account_id, akas, partition, region, tags, title
from aws.aws_lambda_function
where name = '{{ resourceName }}'