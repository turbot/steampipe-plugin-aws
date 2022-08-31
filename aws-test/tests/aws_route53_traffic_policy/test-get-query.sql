select name, id, type, version::varchar
from aws_route53_traffic_policy
where id = '{{ output.id.value }}' and version = '{{ output.version.value }}'