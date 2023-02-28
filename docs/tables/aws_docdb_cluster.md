# Table: aws_docdb_cluster

Amazon DocumentDB (with MongoDB compatibility) is a fast, reliable, and fully managed database service. Amazon DocumentDB makes it easy to set up, operate, and scale MongoDB-compatible databases in the cloud. With Amazon DocumentDB, you can run the same application code and use the same drivers and tools that you use with MongoDB.

## Examples

## Basic Info

```sql
select
  arn,
  db_cluster_identifier,
  deletion_protection,
  engine,
  status,
  region
from
  aws_docdb_cluster;
```

### List clusters which are not encrypted

```sql
select
  db_cluster_identifier,
  status,
  cluster_create_time,
  kms_key_id,
  storage_encrypted
from
  aws_docdb_cluster
where
  not storage_encrypted;
```

### List clusters where backup retention period is greater than 7 days

```sql
select
  db_cluster_identifier,
  backup_retention_period
from
  aws_docdb_cluster
where
  backup_retention_period > 7;
```

### Get avalability zone count for each cluster

```sql
select
  db_cluster_identifier,
  jsonb_array_length(availability_zones) as availability_zones_count
from
  aws_docdb_cluster;
```

### List cluster members details

```sql
select
  db_cluster_identifier,
  member ->> 'DBClusterParameterGroupStatus' as db_cluster_parameter_group_status,
  member ->> 'DBInstanceIdentifier' as db_instance_identifier,
  member ->> 'IsClusterWriter' as is_cluster_writer,
  member ->> 'PromotionTier' as promotion_tier
from
  aws_docdb_cluster
  cross join jsonb_array_elements(members) as member;
```

### List clusters where deletion protection is disabled

```sql
select
  db_cluster_identifier,
  status,
  cluster_create_time,
  deletion_protection
from
  aws_docdb_cluster
where
  not deletion_protection;
```