select
  akas,
  id
from 
  aws_route53_health_check
where title = '{{ output.health_check_id.value }}';