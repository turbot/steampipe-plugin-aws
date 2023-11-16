---
title: "Table: aws_efs_file_system - Query AWS Elastic File System using SQL"
description: "Allows users to query AWS Elastic File System (EFS) file systems, providing detailed information about each file system such as its ID, ARN, creation token, performance mode, and lifecycle state."
---

# Table: aws_efs_file_system - Query AWS Elastic File System using SQL

The `aws_efs_file_system` table in Steampipe provides information about file systems within AWS Elastic File System (EFS). This table allows DevOps engineers to query file system-specific details, including its ID, ARN, creation token, performance mode, lifecycle state, and associated metadata. Users can utilize this table to gather insights on file systems, such as their performance mode, lifecycle state, and more. The schema outlines the various attributes of the EFS file system, including the file system ID, creation token, tags, and associated mount targets.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_efs_file_system` table, you can use the `.inspect aws_efs_file_system` command in Steampipe.

Key columns:

- `file_system_id`: The ID of the file system. This is a unique identifier that can be used to join this table with other tables.
- `owner_id`: The AWS account that created the file system. This can be useful for auditing and managing access control.
- `creation_token`: The opaque string specified in the request to ensure idempotent creation. This can be useful to verify the creation process of the file system.

## Examples

### Basic info

```sql
select
  name,
  file_system_id,
  owner_id,
  automatic_backups,
  creation_token,
  creation_time,
  life_cycle_state,
  number_of_mount_targets,
  performance_mode,
  throughput_mode
from
  aws_efs_file_system;
```


### List file systems which are not encrypted at rest

```sql
select
  file_system_id,
  encrypted,
  kms_key_id,
  region
from
  aws_efs_file_system
where
  not encrypted;
```


### Get the size of the data stored in each file system

```sql
select
  file_system_id,
  size_in_bytes ->> 'Value' as data_size,
  size_in_bytes ->> 'Timestamp' as data_size_timestamp,
  size_in_bytes ->> 'ValueInIA' as data_size_infrequent_access_storage,
  size_in_bytes ->> 'ValueInStandard' as data_size_standard_storage
from
  aws_efs_file_system;
```


### List file systems which have root access

```sql
select
  title,
  p as principal,
  a as action,
  s ->> 'Effect' as effect,
  s -> 'Condition' as conditions
from
  aws_efs_file_system,
  jsonb_array_elements(policy_std -> 'Statement') as s,
  jsonb_array_elements_text(s -> 'Principal' -> 'AWS') as p,
  jsonb_array_elements_text(s -> 'Action') as a
where
  a in ('elasticfilesystem:clientrootaccess');
```


### List file systems that do not enforce encryption in transit

```sql
select
  title
from
  aws_efs_file_system
where
  title not in (
    select
      title
    from
      aws_efs_file_system,
      jsonb_array_elements(policy_std -> 'Statement') as s,
      jsonb_array_elements_text(s -> 'Principal' -> 'AWS') as p,
      jsonb_array_elements_text(s -> 'Action') as a,
      jsonb_array_elements_text(
        s -> 'Condition' -> 'Bool' -> 'aws:securetransport'
      ) as ssl
    where
      p = '*'
      and s ->> 'Effect' = 'Deny'
      and ssl :: bool = false
  );
```


### List file systems with automatic backups enabled

```sql
select
  name,
  automatic_backups,
  arn,
  file_system_id
from
  aws_efs_file_system
where
  automatic_backups = 'enabled';
```
