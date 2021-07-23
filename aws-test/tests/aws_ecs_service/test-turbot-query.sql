select
  title,
  akas
from
  aws.aws_ecs_service
where
  arn = '{{ output.resource_aka.value }}';