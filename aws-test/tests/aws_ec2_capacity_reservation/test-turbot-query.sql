select akas, title
from aws.aws_ec2_capacity_reservation
where capacity_reservation_id = '{{ output.resource_id.value }}';
