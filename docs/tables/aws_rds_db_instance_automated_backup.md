---
title: "Table: aws_rds_db_instance_automated_backup - Query AWS RDS DB Instance Automated Backups using SQL"
description: "Allows users to query AWS RDS DB Instance Automated Backups and retrieve data about automated backups for RDS DB instances."
---

# Table: aws_rds_db_instance_automated_backup - Query AWS RDS DB Instance Automated Backups using SQL

The `aws_rds_db_instance_automated_backup` table in Steampipe allows users to query AWS RDS DB Instance Automated Backups. This table provides data about automated backups for RDS DB instances. It allows DevOps engineers, database administrators, and other technical professionals to query backup-specific details, including backup status, retention period, and associated metadata. Users can utilize this table to gather insights on backups, such as backup statuses, encrypted backups, verification of backup retention periods, and more. The schema outlines the various attributes of the automated backup, including the backup ARN, backup creation date, backup size, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_rds_db_instance_automated_backup` table, you can use the `.inspect aws_rds_db_instance_automated_backup` command in Steampipe.

**Key columns**:

- `dbi_resource_id`: The identifier for the source DB instance, which may not be unique across all instances. This column can be used to join with the `aws_rds_db_instance` table.
- `db_instance_identifier`: The user-provided name of the source DB instance. This column is useful for identifying the specific DB instance associated with the backup.
- `allocated_storage`: The amount of storage allocated for the automated backup. This column is useful for tracking storage utilization and planning capacity.

## Examples

### Basic info

```sql
select
  db_instance_identifier,
  arn,
  status,
  allocated_storage,
  encrypted,
  engine
from
  aws_rds_db_instance_automated_backup;
```

### List DB instance automated backups that are not encrypted

```sql
select
  db_instance_identifier,
  arn,
  status,
  backup_target,
  instance_create_time,
  encrypted,
  engine
from
  aws_rds_db_instance_automated_backup
where
  not encrypted;
```

### List DB instance automated backups that are not authenticated through IAM users and roles

```sql
select
  db_instance_identifier,
  iam_database_authentication_enabled,
  status,
  availability_zone,
  dbi_resource_id
from
  aws_rds_db_instance_automated_backup
where
  not iam_database_authentication_enabled;
```

### Get VPC and subnet info for each DB instance automated backup

```sql
select
  b.arn,
  b.vpc_id,
  v.cidr_block,
  v.is_default,
  v.instance_tenancy
from
  aws_rds_db_instance_automated_backup as b,
  aws_vpc as v
where
  v.vpc_id = b.vpc_id;
```

### List DB instance automated backups of deleted instances

```sql
select
  db_instance_identifier,
  arn,
  engine,
  engine_version,
  availability_zone,
  backup_retention_period,
  status
from
  aws_rds_db_instance_automated_backup
where
  status = 'retained';
```

### Get KMS key details of each DB instance automated backup

```sql
select
  b.db_instance_identifier,
  b.arn as automated_backup_arn,
  b.engine,
  b.kms_key_id,
  k.creation_date as kms_key_creation_date,
  k.key_state,
  k.key_rotation_enabled
from
  aws_rds_db_instance_automated_backup as b,
  aws_kms_key as k
where
  k.id = b.kms_key_id;
```
