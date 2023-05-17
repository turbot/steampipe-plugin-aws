select launch_template_name, launch_template_id, version_number
from aws.aws_ec2_launch_template_version
where launch_template_id = '{{ output.resource_id.value }}' and version_number = 1;
