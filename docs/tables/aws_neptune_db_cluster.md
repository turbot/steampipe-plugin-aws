# Table: aws_neptune_db_cluster

An Amazon Neptune DB cluster manages access to your data through queries.

**Note**: This table only returns Neptune DB clusters, not RDS or DocumentDB DB clusters.

## Examples

### List of DB clusters which are not encrypted

```sql
select
  db_cluster_identifier,
  allocated_storage,
  kms_key_id
from
  aws_neptune_db_cluster
where
  kms_key_id is null;
```

### List of DB clusters where backup retention period is greater than 7 days

```sql
select
  db_cluster_identifier,
  backup_retention_period
from
  aws_neptune_db_cluster
where
  backup_retention_period > 7;
```

### Avalability zone count for each db instance

```sql
select
  db_cluster_identifier,
  jsonb_array_length(availability_zones) availability_zones_count
from
  aws_neptune_db_cluster;
```

### DB cluster Members info

```sql
select
  db_cluster_identifier,
  member ->> 'DBClusterParameterGroupStatus' as db_cluster_parameter_group_status,
  member ->> 'DBInstanceIdentifier' as db_instance_identifier,
  member ->> 'IsClusterWriter' as is_cluster_writer,
  member ->> 'PromotionTier' as promotion_tier
from
  aws_neptune_db_cluster
  cross join jsonb_array_elements(db_cluster_members) as member;
```
