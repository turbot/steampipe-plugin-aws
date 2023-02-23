select image_id, akas, tags, title
from aws.aws_ec2_ami
where image_id = '{{ output.resource_id.value }}'
