select
  service_name,
  arn,
  cluster_arn,
  task_definition,
  desired_count,
  enable_ecs_managed_tags,
  enable_execute_command,
  placement_constraints,
  placement_strategy
from
  aws.aws_ecs_service
where
  arn = '{{ output.resource_aka.value }}';