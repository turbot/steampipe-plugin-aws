# Table: aws_eks_node_group

Amazon EKS managed node groups automate the provisioning and lifecycle management of nodes (Amazon EC2 instances) for Amazon EKS Kubernetes clusters.

## Examples

### Basic info

```sql
select
  nodegroup_name,
  arn,
  created_at,
  cluster_name,
  status
from
  aws_eks_node_group;
```

### List node groups that are not active

```sql
select
  nodegroup_name,
  arn,
  created_at,
  cluster_name,
  status
from
  aws_eks_node_group
where
  status <> 'ACTIVE';
```

### Get health status of the node groups

```sql
select
  nodegroup_name,
  cluster_name,
  jsonb_pretty(health) as health
from
  aws_eks_node_group;
```

### Get launch template details of the node groups

```sql
select
  nodegroup_name,
  cluster_name,
  jsonb_pretty(launch_template) as launch_template
from
  aws_eks_node_group;
```
