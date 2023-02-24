select launch_template_name, launch_template_id
from aws.aws_ec2_launch_template
where launch_template_id = '{{ output.resource_id.value }}'
