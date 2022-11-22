# Table: aws_eks_cluster

The Fargate profile allows an administrator to declare which pods run on Fargate. Each profile can have up to five selectors that contain a namespace and optional labels. You must define a namespace for every selector. The label field consists of multiple optional key-value pairs. Pods that match a selector (by matching a namespace for the selector and all of the labels specified in the selector) are scheduled on Fargate.

## Examples

### Basic info

```sql
select
  fargate_profile_name,
  fargate_profile_arn,
  cluster_name,
  created_at,
  status,
  tags
from
  aws_eks_fargate_profile;
```

### List fargate profiles which are inactive
```sql
select
  fargate_profile_name,
  fargate_profile_arn,
  cluster_name,
  created_at,
  status
from
  aws_eks_fargate_profile
where
  status <> 'ACTIVE';
```

### Get the subnet configuration for each fargate profile

```sql
select
  f.fargate_profile_name as fargate_profile_name,
  f.cluster_name as cluster_name,
  f.status as fargate_profile_status,
  s.availability_zone as availability_zone,
  s.available_ip_address_count as available_ip_address_count,
  s.cidr_block as cidr_block,
  s.vpc_id as vpc_id
from
  aws_eks_fargate_profile as f,
  aws_vpc_subnet as s,
  jsonb_array_elements(f.subnets) as subnet_id
where
  s.subnet_id = subnet_id;
```

### List fargate profiles for clusters not running Kubernetes version greater than 1.19

```sql
select
  c.name as cluster_name,
  c.arn as cluster_arn,
  c.version as cluster_version,
  f.fargate_profile_name as fargate_profile_name,
  f.fargate_profile_arn as fargate_profile_arn,
  f.created_at as created_at,
  f.pod_execution_role_arn as pod_execution_role_arn,
  f.status as fargate_profile_status
from
  aws_eks_fargate_profile as f,
  aws_eks_cluster as c
where
  c.version::float > 1.19 and f.cluster_name = c.name;
```
