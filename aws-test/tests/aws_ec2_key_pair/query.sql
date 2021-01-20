select key_name, key_pair_id
from aws.aws_ec2_key_pair
where key_pair_id = '{{ output.resource_id.value }}'
