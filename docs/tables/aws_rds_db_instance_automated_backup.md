---
title: "Steampipe Table: aws_rds_db_instance_automated_backup - Query AWS RDS DB Instance Automated Backups using SQL"
description: "Allows users to query AWS RDS DB Instance Automated Backups and retrieve data about automated backups for RDS DB instances."
folder: "Backup"
---

# Table: aws_rds_db_instance_automated_backup - Query AWS RDS DB Instance Automated Backups using SQL

The AWS RDS DB Instance Automated Backup is a feature of Amazon RDS that enables automated backups of your DB instances. These backups include transaction logs so that you can perform a point-in-time recovery of your databases. Automated backups are kept for a specified, configurable period, allowing you to restore the database to any point in time during that period.

## Table Usage Guide

The `aws_rds_db_instance_automated_backup` table in Steampipe allows you to query AWS RDS DB Instance Automated Backups. This table provides you with data about automated backups for RDS DB instances. It enables you, as a DevOps engineer, database administrator, or other technical professional, to query backup-specific details, including backup status, retention period, and associated metadata. You can utilize this table to gather insights on backups, such as backup statuses, encrypted backups, verification of backup retention periods, and more. The schema outlines the various attributes of the automated backup for you, including the backup ARN, backup creation date, backup size, and associated tags.

## Examples

### Basic info
Discover the segments that are encrypted within your automated backup instances on AWS RDS, enabling you to assess the elements within your database that are secure. This is particularly useful when managing data security and ensuring compliance with data protection regulations.

```sql+postgres
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

```sql+sqlite
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
Identify instances where automated backups of your database are not encrypted. This can be useful to enhance your data security by ensuring all backups are encrypted.

```sql+postgres
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

```sql+sqlite
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
  encrypted = 0;
```

### List DB instance automated backups that are not authenticated through IAM users and roles
Identify instances where automated backups of the database are not authenticated through IAM users and roles. This is useful for ensuring all backups have the necessary security measures in place.

```sql+postgres
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

```sql+sqlite
select
  db_instance_identifier,
  iam_database_authentication_enabled,
  status,
  availability_zone,
  dbi_resource_id
from
  aws_rds_db_instance_automated_backup
where
  iam_database_authentication_enabled = 0;
```

### Get VPC and subnet info for each DB instance automated backup
This example helps you analyze the relationship between your automated backup instances for your database and their associated virtual private clouds (VPC) and subnets. It's useful for understanding your infrastructure setup and how your database backups are distributed across different VPCs.

```sql+postgres
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

```sql+sqlite
select
  b.arn,
  b.vpc_id,
  v.cidr_block,
  v.is_default,
  v.instance_tenancy
from
  aws_rds_db_instance_automated_backup as b
join
  aws_vpc as v
on
  v.vpc_id = b.vpc_id;
```

### List DB instance automated backups of deleted instances
Discover the segments that are retaining automated backups of deleted database instances. This can be helpful in identifying instances where you may want to free up storage or ensure data from deleted instances is properly archived.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which each database instance's automated backup uses a specific Key Management Service (KMS) key. This can help in understanding the security measures in place and the overall configuration of database backups.

```sql+postgres
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

```sql+sqlite
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