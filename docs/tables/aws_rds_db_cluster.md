---
title: "Table: aws_rds_db_cluster - Query AWS RDS DB Clusters using SQL"
description: "Allows users to query AWS RDS DB Clusters and retrieve valuable information about the status, configuration, and security settings of each DB cluster."
---

# Table: aws_rds_db_cluster - Query AWS RDS DB Clusters using SQL

The `aws_rds_db_cluster` table in Steampipe provides information about DB clusters within Amazon Relational Database Service (RDS). This table allows DevOps engineers to query DB cluster-specific details, including configuration, status, and security settings. Users can utilize this table to gather insights on DB clusters, such as their availability, backup settings, encryption status, and more. The schema outlines the various attributes of the DB cluster, including the DB cluster identifier, creation time, DB cluster members, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_rds_db_cluster` table, you can use the `.inspect aws_rds_db_cluster` command in Steampipe.

**Key columns**:

- `db_cluster_identifier`: The identifier for the DB cluster. This identifier is unique in the scope of an AWS account and can be used to join this table with other AWS RDS tables.
- `status`: The status of the DB cluster, such as 'available', 'modifying', etc. This is important for monitoring the availability and health of the DB cluster.
- `creation_time`: The timestamp when the DB cluster was created. This can be useful for auditing and tracking the lifecycle of DB clusters.

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
