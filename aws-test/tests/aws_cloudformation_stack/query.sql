select name, tags_raw
from aws.aws_cloudformation_stack
where name = '{{resourceName}}'
order by name