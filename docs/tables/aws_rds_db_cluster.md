# Table: aws_rds_db_cluster

An Amazon Aurora DB cluster consists of one or more DB instances and a cluster volume that manages the data for those DB instances.

**Note**: This table only returns RDS DB clusters, e.g., Aurora, MySQL, Postgres, not DocumentDB or Neptune DB clusters.

## Examples

### List of DB clusters which are not encrypted

```sql
select
  db_cluster_identifier,
  allocated_storage,
  kms_key_id
from
  aws_rds_db_cluster
where
  kms_key_id is null;
```

### List of DB clusters where backup retention period is greater than 7 days

```sql
select
  db_cluster_identifier,
  backup_retention_period
from
  aws_rds_db_cluster
where
  backup_retention_period > 7;
```

### Avalability zone count for each db instance

```sql
select
  db_cluster_identifier,
  jsonb_array_length(availability_zones) availability_zones_count
from
  aws_rds_db_cluster;
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
  aws_rds_db_cluster
  cross join jsonb_array_elements(members) as member;
```

### List DB cluster pending maintenance actions

```sql
select
  actions ->> 'ResourceIdentifier' as db_cluster_identifier,
  details ->> 'Action' as action,
  details ->> 'OptInStatus' as opt_in_status,
  details ->> 'ForcedApplyDate' as forced_apply_date,
  details ->> 'CurrentApplyDate' as current_apply_date,
  details ->> 'AutoAppliedAfterDate' as auto_applied_after_date
from
  aws_rds_db_cluster,
  jsonb_array_elements(pending_maintenance_actions) as actions,
  jsonb_array_elements(actions -> 'PendingMaintenanceActionDetails') as details;
```
