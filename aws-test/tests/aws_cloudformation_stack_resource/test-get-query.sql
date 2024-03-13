select logical_resource_id, stack_name, stack_id
from aws.aws_cloudformation_stack_resource
where stack_name = '{{ resourceName }}' and logical_resource_id = 'myVpc';
