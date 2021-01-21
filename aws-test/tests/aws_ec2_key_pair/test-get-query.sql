select key_name, key_pair_id, key_fingerprint, tags_src
from aws.aws_ec2_key_pair
where key_name = '{{resourceName}}'
