select id, accepter_vpc_id, requester_vpc_id
from aws.aws_vpc_peering_connection
where id = '{{ output.id.value }}';
