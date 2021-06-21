select akas, title
from aws.aws_ec2_ami_shared
where image_id = '{{ output.resource_id.value }}';
