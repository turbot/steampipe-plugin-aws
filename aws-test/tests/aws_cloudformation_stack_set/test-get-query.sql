select stack_set_name, stack_set_id, status, title
from aws_cloudformation_stack_set
where stack_set_name = '{{resourceName}}'
