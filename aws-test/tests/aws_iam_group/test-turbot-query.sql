select title, akas, account_id, partition, region
from aws.aws_iam_group
where name = '{{ resourceName }}'
