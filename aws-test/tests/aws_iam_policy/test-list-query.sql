select name, arn
from aws.aws_iam_policy
where arn = '{{ output.resource_aka.value }}'
