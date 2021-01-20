select name, arn
from aws.aws_iam_user
where arn = '{{ output.resource_aka.value }}'
