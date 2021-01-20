select title, akas, region, partition, account_id
from aws.aws_iam_access_key
where access_key_id = '{{ output.resource_id.value }}'
