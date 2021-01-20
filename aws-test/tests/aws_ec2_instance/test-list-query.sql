select instance_id, image_id
from aws.aws_ec2_instance
where instance_id = '{{ output.resource_id.value }}'
