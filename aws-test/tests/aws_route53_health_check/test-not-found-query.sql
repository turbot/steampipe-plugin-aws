select
  akas,
  id,
  tags_src
from 
  aws_route53_health_check
where id = 'dummy{{ output.health_check_id.value }}';