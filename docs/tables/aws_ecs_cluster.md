# Table: aws_ecs_cluster

An Amazon ECS cluster is a logical grouping of tasks or services. If you are running tasks or services that use the EC2 launch type, a cluster is also a grouping of container instances. If you are using capacity providers, a cluster is also a logical grouping of capacity providers.

## Examples

### Basic info

```sql
select
  cluster_arn,
  cluster_name,
  active_services_count,
  attachments,
  attachments_status,
  status
from
  aws_ecs_cluster;
```


### List clusters that have failed to provision resources

```sql
select
  cluster_arn,
  status
from
  aws_ecs_cluster
where
  status = 'FAILED';
```


### Get details of resources attached to each cluster

```sql
select
  cluster_arn,
  attachment ->> 'id' as attachment_id,
  attachment ->> 'status' as attachment_status,
  attachment ->> 'type' as attachment_type
from
  aws_ecs_cluster,
  jsonb_array_elements(attachments) as attachment;
```


### List clusters with CloudWatch Container Insights disabled

```sql
select
  cluster_arn,
  setting ->> 'Name' as name,
  setting ->> 'Value' as value
from
  aws_ecs_cluster,
  jsonb_array_elements(settings) as setting
where
  setting ->> 'Value' = 'disabled';
```
