# Table: aws_ecs_cluster

An Amazon ECS cluster is a logical grouping of tasks or services. If you are running tasks or services that use the EC2 launch type, a cluster is also a grouping of container instances. If you are using capacity providers, a cluster is also a logical grouping of capacity providers.

## Examples

### Basic ecs cluster info

```sql
select
  cluster_arn,
  cluster_name,
  active_sevices_count,
  attachments,
  attachments_status,
  capacity_providers,
  default_capacity_provider_strategy,
  pending_tasks_count,
  registered_container_instances_count,
  running_tasks_count,
  settings,
  statistics,
  status
from
  aws_new.aws_ecs_cluster;
```


### List of ecs clusters with state in FAILED or INACTIVE

```sql
select
  cluster_arn,
  status
from
  aws_new.aws_ecs_cluster
where
  status IN ('FAILED','INACTIVE');
```


### List resources attached to ecs clusters

```sql
select
  attachments,
  attachments_status
from
  aws_new.aws_ecs_cluster;
```


### List of cluster for which CloudWatch Container Insights is disabled

```sql
select
  cluster_arn,
  setting ->> 'Value' as containerInsights
from
  aws_new.aws_ecs_cluster,
  jsonb_array_elements(settings) as setting;
```
