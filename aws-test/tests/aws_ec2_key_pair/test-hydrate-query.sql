select key_name, key_pair_id, akas, tags, title
from aws.aws_ec2_key_pair
where key_name = '{{resourceName}}'
