# Table: aws_dax_parameter_group

Amazon DynamoDB Accelerator (DAX) Parameter Group is a named set of parameters that are applied to every node in the cluster. You can use a parameter group to specify cache TTL behavior.

## Examples

### Basic info

```sql
select
  parameter_group_name,
  description,
  region
from
  aws_dax_parameter_group;
```

### Get cluster details that associated with parameter group

```sql
select
  p.parameter_group_name,
  c.cluster_name,
  c.node_type,
  c.status
from
  aws_dax_parameter_group as p,
  aws_dax_cluster as c
where
  c.parameter_group ->> 'ParameterGroupName' = p.parameter_group_name;
```