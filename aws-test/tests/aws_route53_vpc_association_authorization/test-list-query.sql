select hosted_zone_id, akas
from aws_route53_vpc_association_authorization
where hosted_zone_id = '{{output.resource_hosted_zone_id.value}}'
and vpc_id = '{{output.resource_vpc_alternate_id.value}}'
