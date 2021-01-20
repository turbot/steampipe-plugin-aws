select name, id, tags, title, akas
from aws.aws_cloudformation_stack
where name = '{{resourceName}}'