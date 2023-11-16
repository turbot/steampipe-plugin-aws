---
title: "Table: aws_neptune_db_cluster - Query Amazon Neptune DB Clusters using SQL"
description: "Allows users to query Amazon Neptune DB clusters for comprehensive information about their configuration, status, and other relevant details."
---

# Table: aws_neptune_db_cluster - Query Amazon Neptune DB Clusters using SQL

The `aws_neptune_db_cluster` table in Steampipe provides information about DB clusters within Amazon Neptune. This table allows DevOps engineers to query DB cluster-specific details, including configuration, status, and associated metadata. Users can utilize this table to gather insights on DB clusters, such as their availability, security settings, backup policies, and more. The schema outlines the various attributes of the DB cluster, including the cluster identifier, creation time, enabled cloudwatch logs exports, and associated tags.

**Note**: This table only returns Neptune DB clusters, not RDS or DocumentDB DB clusters.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_neptune_db_cluster` table, you can use the `.inspect aws_neptune_db_cluster` command in Steampipe.

### Key columns:

- `db_cluster_identifier`: The identifier for the DB cluster. This identifier is unique across all DB clusters within an AWS account, and can be used to join this table with other tables.
- `arn`: The Amazon Resource Name (ARN) for the DB cluster. ARN is a globally unique identifier, which can be used to join with other tables where ARN is a common column.
- `vpc_security_groups`: The list of VPC security groups that are associated with the DB cluster. This information can be useful for security audits and compliance checks.

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
