# Table: aws_dax_cluster

Amazon DynamoDB Accelerator (DAX) is a fully managed, highly available, in-memory cache for Amazon DynamoDB that delivers up to a 10 times performance improvement—from milliseconds to microseconds—even at millions of requests per second.

## Examples

### Basic info

```sql
select
  cluster_name,
  description,
  active_nodes,
  iam_role_arn,
  status,
  region
from
  aws_dax_cluster;
```


### List clusters that does not enforce server-side encryption (SSE)

```sql
select
  cluster_name,
  description,
  sse_description ->> 'Status' as sse_status
from
  aws_dax_cluster
where
  sse_description ->> 'Status' = 'DISABLED';
```


### List clusters provisioned with undesired (for example, dax.r4.large and dax.t2.medium are desired) node types

```sql
select
  cluster_name,
  node_type,
  count(*) as count
from
  aws_dax_cluster
where
  node_type not in ('cache.m5.large', 'cache.m4.4xlarge')
group by
  cluster_name, node_type;
```


### Get the network details for each cluster

```sql
select
  cluster_name,
  subnet_group,
  sg ->> 'SecurityGroupIdentifier' as sg_id,
  n ->> 'AvailabilityZone' as az_name,
  cluster_discovery_endpoint ->> 'Address' as cluster_discovery_endpoint_address,
  cluster_discovery_endpoint ->> 'Port' as cluster_discovery_endpoint_port
from
  aws_dax_cluster,
  jsonb_array_elements(security_groups) as sg,
  jsonb_array_elements(nodes) as n;
```