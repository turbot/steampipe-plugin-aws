# Table: aws_msk_serverless_cluster

MSK Serverless is a cluster type for Amazon MSK that makes it possible for you to run Apache Kafka without having to manage and scale cluster capacity. It automatically provisions and scales capacity while managing the partitions in your topic, so you can stream data without thinking about right-sizing or scaling clusters.

## Examples

### Basic Info

```sql
select
  arn,
  cluster_name,
  state,
  cluster_type,
  creation_time,
  current_version,
  region,
  tags
from
  aws_msk_serverless_cluster;
```

### List inactive clusters

```sql
select
  arn,
  cluster_name,
  state,
  creation_time
from
  aws_msk_serverless_cluster
where
  state <> 'ACTIVE';
```

### List clusters created within the last 90 days

```sql
select
  arn,
  cluster_name,
  state,
  creation_time
from
  aws_msk_serverless_cluster
where
  creation_time >= (current_date - interval '90' day)
order by
  creation_time;
```

### Get VPC details of each cluster

```sql
select
  arn,
  cluster_name,
  state,
  vpc ->> 'SubnetIds' as subnet_ids,
  vpc ->> 'SecurityGroupIds' as security_group_ids
from
  aws_msk_serverless_cluster,
  jsonb_array_elements(serverless -> 'VpcConfigs') as vpc
```

### List clusters with IAM authentication disabled

```sql
select
  arn,
  cluster_name,
  state,
  serverless -> 'ClientAuthentication' as client_authentication
from
  aws_msk_serverless_cluster
where
  (serverless -> 'ClientAuthentication' -> 'Sasl' -> 'Iam' ->> 'Enabled')::boolean = false
```
