# Table: aws_emr_instance

The central component of Amazon EMR is the cluster. A cluster is a collection of Amazon Elastic Compute Cloud (Amazon EC2) instances. Each instance in the cluster is called a node. Each node has a role within the cluster, referred to as the node type.

## Examples

### Basic info

```sql
select
  id,
  cluster_id,
  ec2_instance_id,
  instance_type,
  private_dns_name,
  private_ip_address
from
  aws_emr_instance;
```

### List instances by type

```sql
select
  id,
  ec2_instance_id,
  instance_type
from
  aws_emr_instance
where
  instance_type = 'm2.4xlarge';
```

### List instances for a cluster

```sql
select
  id,
  ec2_instance_id,
  instance_type
from
  aws_emr_instance
where
  cluster_id = 'j-21HIX5R2NZMXJ';
```

### Get volume details for an instance

```sql
select
  id,
  ec2_instance_id,
  instance_type,
  v -> 'Device' as device,
  v -> 'VolumeId' as volume_id
from
  aws_emr_instance,
  jsonb_array_elements(ebs_volumes) as v
where
  id = 'ci-ULCFS2ZN0FK7';
```