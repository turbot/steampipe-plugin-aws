select name, id, comment
from aws_route53_zone
where akas::text = '["{{ output.resource_aka.value }}"]';