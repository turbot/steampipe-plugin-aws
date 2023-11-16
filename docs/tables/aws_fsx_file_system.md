---
title: "Table: aws_fsx_file_system - Query AWS FSx File Systems using SQL"
description: "Allows users to query AWS FSx File Systems to gather information about the file system's details, including its lifecycle, type, storage capacity, and associated tags."
---

# Table: aws_fsx_file_system - Query AWS FSx File Systems using SQL

The `aws_fsx_file_system` table in Steampipe provides information about FSx File Systems within Amazon Web Services. This table allows DevOps engineers to query file system-specific details, including its lifecycle, type, storage capacity, and associated tags. Users can utilize this table to gather insights on file systems, such as their storage capacity, creation times, and statuses. The schema outlines the various attributes of the FSx File Systems, including the file system ID, owner ID, creation time, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_fsx_file_system` table, you can use the `.inspect aws_fsx_file_system` command in Steampipe.

### Key columns:

- `file_system_id`: This is the unique identifier of the FSx File System. It can be used to join this table with other tables that contain information about AWS resources.
- `owner_id`: This is the AWS account ID of the file system owner. It can be used to filter file systems by owner.
- `creation_time`: This is the timestamp of when the file system was created. It can be used to understand the lifecycle and age of the file system.

## Examples

### Basic info

```sql
select
  file_system_id,
  arn,
  dns_name,
  owner_id,
  creation_time,
  lifecycle,
  storage_capacity
from
  aws_fsx_file_system;
```

### List file systems which are encrypted

```sql
select
  file_system_id,
  kms_key_id,
  region
from
  aws_fsx_file_system
where
  kms_key_id is not null;
```