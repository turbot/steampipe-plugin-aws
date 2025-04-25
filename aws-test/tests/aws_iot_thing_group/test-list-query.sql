select group_name, arn, version::text, parent_group_name
from aws.aws_iot_thing_group
where parent_group_name = '{{ output.parent_resource_name.value }}';