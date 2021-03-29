select task_definition_arn, cpu, family, memory, status, account_id, region, partition
from aws.aws_ecs_task_definition
where task_definition_arn = '{{ output.resource_aka.value }}';