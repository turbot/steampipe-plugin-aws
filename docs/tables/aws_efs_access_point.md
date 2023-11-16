---
title: "Table: aws_efs_access_point - Query Amazon EFS Access Points using SQL"
description: "Allows users to query Amazon EFS Access Points, providing detailed information about each access point's configuration, including the file system it is associated with, its access point ID, and other related metadata."
---

# Table: aws_efs_access_point - Query Amazon EFS Access Points using SQL

The `aws_efs_access_point` table in Steampipe provides information about Access Points within Amazon Elastic File System (EFS). This table allows DevOps engineers, system administrators, and other technical professionals to query access point-specific details, including the file system it is associated with, its access point ID, and other related metadata. Users can utilize this table to gather insights on access points, such as their operating system type, root directory creation info, and more. The schema outlines the various attributes of the access point, including the access point ARN, creation time, life cycle state, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_efs_access_point` table, you can use the `.inspect aws_efs_access_point` command in Steampipe.

### Key columns:

- `access_point_id`: The unique identifier of the access point. This column can be used to join with other tables when you need to retrieve specific information about a particular access point.
- `file_system_id`: The ID of the EFS file system that the access point applies to. This column is useful when you need to correlate access points with their associated file systems.
- `access_point_arn`: The Amazon Resource Name (ARN) associated with the access point. This column is useful when you need to join with other AWS resource tables using ARN.

## Examples

### Basic info

```sql
select
  name,
  access_point_id,
  access_point_arn,
  file_system_id,
  life_cycle_state,
  owner_id,
  root_directory
from
  aws_efs_access_point;
```


### List access points for each file system

```sql
select
  name,
  access_point_id,
  file_system_id,
  owner_id,
  root_directory
from
  aws_efs_access_point
```


### List access points in the error lifecycle state

```sql
select
  name,
  access_point_id,
  life_cycle_state,
  file_system_id,
  owner_id,
  root_directory
from
  aws_efs_access_point
where
  life_cycle_state = 'error';
```
