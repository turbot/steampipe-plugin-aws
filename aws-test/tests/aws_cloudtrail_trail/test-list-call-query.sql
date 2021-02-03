select arn, name
from aws.aws_cloudtrail_trail
where arn = '{{ output.resource_aka.value }}'
