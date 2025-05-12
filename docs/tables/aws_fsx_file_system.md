---
title: "Steampipe Table: aws_fsx_file_system - Query AWS FSx File Systems using SQL"
description: "Allows users to query AWS FSx File Systems to gather information about the file system's details, including its lifecycle, type, storage capacity, and associated tags."
folder: "FSx"
---

# Table: aws_fsx_file_system - Query AWS FSx File Systems using SQL

The AWS FSx File System is a fully managed service that makes it easy to launch and run feature-rich and highly-performant file systems. With FSx, you can leverage the rich feature sets and fast performance of widely-used open source and commercially-licensed file systems, while avoiding time-consuming administrative tasks. FSx is built on Windows Server, delivering a wide range of administrative features such as user quotas, end-user file restore, and Microsoft Active Directory integration.

## Table Usage Guide

The `aws_fsx_file_system` table in Steampipe provides you with information about FSx File Systems within Amazon Web Services. This table allows you, as a DevOps engineer, to query file system-specific details, including its lifecycle, type, storage capacity, and associated tags. You can utilize this table to gather insights on file systems, such as their storage capacity, creation times, and statuses. The schema outlines the various attributes of the FSx File Systems for you, including the file system ID, owner ID, creation time, and associated tags.

## Examples

### Basic info
Explore the fundamental characteristics of your AWS FSx file systems, such as ownership details, creation time, and storage capacity. This can help manage resources and identify potential areas for optimization or troubleshooting.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that are encrypted within your file systems to enhance your data security and privacy measures. This query is useful in identifying those file systems that have encryption enabled, ensuring compliance with data protection regulations.

```sql+postgres
select
  file_system_id,
  kms_key_id,
  region
from
  aws_fsx_file_system
where
  kms_key_id is not null;
```

```sql+sqlite
select
  file_system_id,
  kms_key_id,
  region
from
  aws_fsx_file_system
where
  kms_key_id is not null;
```