select name, arn
from aws.aws_iam_role
where arn = '{{ output.resource_aka.value }}'
