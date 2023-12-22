select launch_template_name, version_number
from aws.aws_ec2_launch_template_version
where launch_template_name = '{{ resourceName }}';
