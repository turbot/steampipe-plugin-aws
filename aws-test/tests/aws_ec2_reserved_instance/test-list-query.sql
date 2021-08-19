select reserved_instance_id, arn, instance_type
from aws.aws_ec2_reserved_instance
where instance_type = '{{ output.instance_type.value }}';
