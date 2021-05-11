# Table: aws_ecs_container_instance

An Amazon ECS container instance is a EC2 instance that is running the Amazon ECS agent and has been registered with an
ECS cluster.

## Examples

### Basic info

```sql
select
  arn,
  ec2_instance_id,
  status,
  status_reason,
  running_tasks_count,
  pending_tasks_count
from
  aws_ecs_container_instance;
```


### List container instances that have failed registration

```sql
select
  arn,
  status,
  status_reason
from
  aws_ecs_container_instance
where
  status = 'REGISTRATION_FAILED';
```


### Get details of resources attached to each container instance

```sql
select
  arn,
  attachment ->> 'id' as attachment_id,
  attachment ->> 'status' as attachment_status,
  attachment ->> 'type' as attachment_type
from
  aws_ecs_container_instance,
  jsonb_array_elements(attachments) as attachment;
```


### List container instances with using a given AMI

```sql
select
  arn,
  setting ->> 'Name' as name,
  setting ->> 'Value' as value
from
  aws_ecs_container_instance,
  jsonb_array_elements(attributes) as setting
where
  setting ->> 'Name' = 'ecs.ami-id' and
  setting ->> 'Value' = 'ami-0babb0c4a4e5769b8';
```
