---
title: "Table: aws_dms_replication_instance - Query AWS Database Migration Service Replication Instances using SQL"
description: "Allows users to query AWS Database Migration Service Replication Instances and provides information about each replication instance in an AWS DMS (Database Migration Service)."
---

# Table: aws_dms_replication_instance - Query AWS Database Migration Service Replication Instances using SQL

The `aws_dms_replication_instance` table in Steampipe provides information about each replication instance in an AWS Database Migration Service. This table allows database administrators to query replication-specific details, including engine version, instance class, allocated storage, and associated metadata. Users can utilize this table to gather insights on replication instances, such as their current state, multi-AZ mode, publicly accessible status, and more. The schema outlines the various attributes of the replication instance, including the replication instance ARN, replication instance identifier, availability zone, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_dms_replication_instance` table, you can use the `.inspect aws_dms_replication_instance` command in Steampipe.

Key columns:

- `replication_instance_arn`: The Amazon Resource Name (ARN) of the replication instance. This is a unique identifier and can be used to join this table with other tables that reference the replication instance ARN.
- `replication_instance_identifier`: The user-defined replication instance identifier. This column can be useful to join with custom tables that use the replication instance identifier as a reference.
- `availability_zone`: The name of the availability zone where the replication instance is located. This column can be used to join this table with other tables that reference the availability zone.

## Examples

### Basic info

```sql
select
  replication_instance_identifier,
  arn,
  engine_version,
  instance_create_time,
  kms_key_id,
  publicly_accessible,
  region
from
  aws_dms_replication_instance;
```


### List replication instances with auto minor version upgrades disabled

```sql
select
  replication_instance_identifier,
  arn,
  engine_version,
  instance_create_time,
  auto_minor_version_upgrade,
  region
from
  aws_dms_replication_instance
where
  not auto_minor_version_upgrade;
```


### List replication instances provisioned with undesired (for example, dms.r5.16xlarge and dms.r5.24xlarge are not desired) instance classes

```sql
select
  replication_instance_identifier,
  arn,
  engine_version,
  instance_create_time,
  replication_instance_class,
  region
from
  aws_dms_replication_instance
where
  replication_instance_class not in ('dms.r5.16xlarge', 'dms.r5.24xlarge');
```


### List publicly accessible replication instances

```sql
select
  replication_instance_identifier,
  arn,
  publicly_accessible,
  region
from
  aws_dms_replication_instance
where
  publicly_accessible;
```


### List replication instances not using multi-AZ deployment configurations

```sql
select
  replication_instance_identifier,
  arn,
  publicly_accessible,
  multi_az,
  region
from
  aws_dms_replication_instance
where
  not multi_az;
```
