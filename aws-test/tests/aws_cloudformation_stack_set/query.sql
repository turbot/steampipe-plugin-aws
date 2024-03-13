select arn, stack_set_id, stack_set_name
from aws_cloudformation_stack_set
where stack_set_name = '{{resourceName}}';
