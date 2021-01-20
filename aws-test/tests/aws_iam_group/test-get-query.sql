select account_id, arn, group_id, name, partition, path
from aws.aws_iam_group
where name = '{{ resourceName }}'
