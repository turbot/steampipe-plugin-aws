select launch_template_name, tags
from aws.aws_ec2_launch_template
where launch_template_name = '{{ resourceName }}'
