select name, tags_src
from aws.aws_cloudformation_stack
where name = '{{resourceName}}'
order by name
