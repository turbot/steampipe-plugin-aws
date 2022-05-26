select
  akas,
  id,
  tags
from 
  aws_route53_health_check
where name = '{{ resourceName }}';