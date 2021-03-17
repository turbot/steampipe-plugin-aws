# Table: aws_eks_node_group

A node group is a collection of managed nodes. Managed nodes are WebSphere® Application Server nodes. A node group defines a boundary for server cluster formation.

## Examples

### Basic info

```sql
select
  name,
  node_group_arn,
  cluster_name,
  version,
  release_version,
  status,
  region
from
  aws_eks_node_group;
```


### List the subnets associated with node group

```sql
select
  name,
  jsonb_array_elements_text(subnets) as subnets
from
  aws_eks_node_group;
```


### List of remote access configuration that is associated with the node group

```sql
select
  name,
  remote_access ->> 'Ec2SshKey' as ec2_ssh_key,
  remote_access ->> 'SourceSecurityGroups' as source_security_groups
from
  aws_eks_node_group;
```


### List the scaling info of node group

```sql
select
  name,
  scaling_config -> 'DesiredSize' as desired_size,
  scaling_config -> 'MaxSize' as max_size,
  scaling_config -> 'MinSize' as min_size
from
  aws_eks_node_group;
```


### List the health info of node groups

```sql
select
  name,
  h ->> 'Code' as code,
  h ->> 'Message' as max_size,
  h -> 'ResourceIds' as resource_ids
from
  aws_eks_node_group,
  jsonb_array_elements(health -> 'Issues') as h;
```