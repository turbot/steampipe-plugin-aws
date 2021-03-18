# Table: aws_ecs_task_definition

A task definition is required to run Docker containers in Amazon ECS.

## Examples

### Basic info

```sql
select
  task_definition_arn,
  cpu,
  title,
  status,
  tags
from
  aws_ecs_task_definition;
```


### Number of containers attached to task definitions

```sql
select
  task_definition_arn,
  jsonb_array_length(container_definitions) as no_of_conatiners
from
  aws_ecs_task_definition;
```


### List of containers with elevated privileges on the host container instance

```sql
select
  task_definition_arn,
  cd ->> 'Privileged' as privileged,
  cd ->> 'Name' as container_name
from
  aws_ecs_task_definition,
  jsonb_array_elements(container_definitions) as cd
where
  cd ->> 'Privileged' = 'true';
```