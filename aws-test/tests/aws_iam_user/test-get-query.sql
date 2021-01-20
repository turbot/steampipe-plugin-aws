select name, path, arn, groups, partition, title, akas, partition, account_id
from aws.aws_iam_user
where name = '{{resourceName}}'
