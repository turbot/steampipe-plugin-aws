select task_definition_arn, title, tags, akas
from aws.aws_ecs_task_definition
where akas::text = '["{{ output.resource_aka.value }}"]';