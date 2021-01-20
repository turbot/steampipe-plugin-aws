select image_id, launch_permissions, akas, tags, title
from aws.aws_ec2_ami
where image_id = '{{ output.resource_id.value }}'
