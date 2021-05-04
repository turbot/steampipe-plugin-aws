select task_definition_arn, cpu, family, memory, status
from aws.aws_ecs_task_definition
where akas::text = '["{{ output.resource_aka.value }}"]';