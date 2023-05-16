select stack_name, logical_resource_id
from aws.aws_cloudformation_stack_resource
where stack_name = '{{ resourceName }}'
