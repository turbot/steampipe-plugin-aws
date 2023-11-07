select stack_set_name, stack_set_id
from aws_cloudformation_stack_set
where arn = '{{ output.resource_aka.value }}'
