select
  akas,
  title,
  tags
from 
  aws_route53_health_check
where id = '{{ output.health_check_id.value }}';