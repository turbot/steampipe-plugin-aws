select name, group_id
from aws.aws_iam_group
where arn = '{{ output.resource_aka.value }}'
