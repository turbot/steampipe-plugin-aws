# Table: aws_msk_cluster

Amazon Managed Streaming for Apache Kafka (Amazon MSK) is a fully managed service that enables you to build and run applications that use Apache Kafka to process streaming data. Amazon MSK provides the control-plane operations, such as those for creating, updating, and deleting clusters. It lets you use Apache Kafka data-plane operations, such as those for producing and consuming data. It runs open-source versions of Apache Kafka.

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
  aws_msk_cluster;
```

### List inactive clusters

```sql
select
  arn,
  cluster_name,
  state,
  creation_time
from
  aws_msk_cluster
where
  state <> 'ACTIVE';
```

### List clusters that allow public access

```sql
select
  arn,
  cluster_name,
  state,
  creation_time
from
  aws_msk_cluster
where
  provisioned -> 'BrokerNodeGroupInfo' -> 'ConnectivityInfo' -> 'PublicAccess' ->> 'Type' <> 'DISABLED';
```

### List clusters with encryption at rest disabled

```sql
select
  arn,
  cluster_name,
  state,
  creation_time
from
  aws_msk_cluster
where
  provisioned -> 'EncryptionInfo' -> 'EncryptionAtRest' is null;
```

### List clusters with encryption in transit disabled

```sql
select
  arn,
  cluster_name,
  state,
  creation_time
from
  aws_msk_cluster
where
  provisioned -> 'EncryptionInfo' -> 'EncryptionInTransit' is null;
```

### List clusters with logging disabled

```sql
select
  arn,
  cluster_name,
  state,
  creation_time
from
  aws_msk_cluster
where
  provisioned -> 'LoggingInfo' is null;
```

### Get total storage used by all the clusters

```sql
select
  sum((provisioned -> 'BrokerNodeGroupInfo' -> 'StorageInfo' -> 'EbsStorageInfo' ->> 'VolumeSize')::int) as total_storage
from
  aws_msk_cluster;
```
