# Table: aws_dms_replication_instance

AWS Database Migration Service (DMS) replication instance is used to connect to your source data store, read the source data, and format the data for consumption by the target data store. A replication instance also loads the data into the target data store.

## Examples

### Basic info

```sql
select
  replication_instance_identifier,
  engine_version,
  instance_create_time,
  kms_key_id,
  publicly_accessible,
  replication_instance_arn,
  region
from
  aws_dms_replication_instance;
```

## TODO
### List clusters that does not enforce server-side encryption (SSE)

```sql
select
  cluster_name,
  description,
  sse_description ->> 'Status' as sse_status
from
  aws_dms_replication_instance
where
  sse_description ->> 'Status' = 'DISABLED';
```

## TODO
### List clusters provisioned with undesired (for example, cache.m5.large and cache.m4.4xlarge are desired) node types

```sql
select
  cluster_name,
  node_type,
  count(*) as count
from
  aws_dms_replication_instance
where
  node_type not in ('cache.m5.large', 'cache.m4.4xlarge')
group by
  cluster_name, node_type;
```


### List replication instances that are publicly accessible

```sql
select
  replication_instance_identifier,
  replication_instance_arn,
  publicly_access,
  region
from
  aws_dms_replication_instance
where
  publicly_access = 'true';
```


### List replication instances not using Multi-AZ deployment configurations

```sql
select
  replication_instance_identifier,
  replication_instance_arn,
  publicly_access,
  multi_az,
  region
from
  aws_dms_replication_instance
where
  multi_az = 'false';
```
