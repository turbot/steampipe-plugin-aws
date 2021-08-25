select capacity_reservation_id, instance_type, state, availability_zone, instance_platform, partition, region, account_id
from aws.aws_ec2_capacity_reservation
where capacity_reservation_arn = '{{ output.resource_aka.value }}';
