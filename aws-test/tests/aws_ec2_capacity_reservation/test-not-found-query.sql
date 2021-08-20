select capacity_reservation_id, instance_type, state
from aws.aws_ec2_capacity_reservation
where capacity_reservation_id = 'dummy-test-{{ output.resource_id.value }}';
