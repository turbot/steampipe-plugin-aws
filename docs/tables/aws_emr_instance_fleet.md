# Table: aws_emr_instance_fleet

The instance fleet configuration for Amazon EMR clusters lets you select a wide variety of provisioning options for Amazon EC2 instances, and helps you develop a flexible and elastic resourcing strategy for each node type in your cluster. You can have only one instance fleet per master, core, and task node type. In an instance fleet configuration, you specify a target capacity for On-Demand Instances and Spot Instances within each fleet.

## Examples

### Basic info

```sql
select
  id,
  arn,
  cluster_id,
  instance_fleet_type,
  state
from
  aws_emr_instance_fleet;
```

### Get the cluster details of the instance fleets

```sql
select
  cluster_id,
  c.name as cluster_name,
  c.state as cluster_state
from
  aws_emr_instance_fleet as f,
  aws_emr_cluster as c
where
  f.cluster_id = c.id;
```

### Get the provisioned & target on demand capacity of the instance fleets

```sql
select
  cluster_id,
  provisioned_on_demand_capacity,
  target_on_demand_capacity
from
  aws_emr_instance_fleet;
```
