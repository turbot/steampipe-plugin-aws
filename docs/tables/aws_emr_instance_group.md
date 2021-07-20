# Table: aws_emr_instance_group

An AWS EMR instance group is a group of instances that have common purpose. With the instance groups configuration, each node type (master, core, or task) consists of the same instance type and the same purchasing option for instances.

## Examples

### Basic info

```sql
select
  id,
  arn,
  cluster_id,
  instance_group_type,
  state
from
  aws_emr_instance_group;
```

### Get the master instance type used for a cluster

```sql
select
  ig.id as instance_group_id,
  ig.cluster_id,
  c.name as cluster_name,
  ig.instance_type
from
  aws_emr_instance_group as ig,
  aws_emr_cluster as c
where
  ig.cluster_id = c.id
  and ig.instance_group_type = 'MASTER';
```

### Get the count of running instances (core and master) per cluster

```sql
select
  cluster_id,
  sum(running_instance_count) as running_instance_count
from
  aws_emr_instance_group
where
  state = 'RUNNING'
group by cluster_id;
```
