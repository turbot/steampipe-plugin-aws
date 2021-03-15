select task_definition_arn, cpu, family, memory, status
from aws_new.aws_ecs_task_definition
where task_definition_arn = '{{ output.resource_aka.value }}';