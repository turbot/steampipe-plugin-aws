select arn, name
from aws.aws_iam_role
where name = '{{resourceName}}'
