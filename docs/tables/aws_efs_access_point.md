---
title: "Steampipe Table: aws_efs_access_point - Query Amazon EFS Access Points using SQL"
description: "Allows users to query Amazon EFS Access Points, providing detailed information about each access point's configuration, including the file system it is associated with, its access point ID, and other related metadata."
folder: "EFS"
---

# Table: aws_efs_access_point - Query Amazon EFS Access Points using SQL

The Amazon Elastic File System (EFS) Access Points provide a customized view into an EFS file system. They enable applications to use a specific operating system user and group, and a directory in the file system as a root directory. By using EFS Access Points, you can enforce a user identity, permission strategy, and root directory for each application using the file system.

## Table Usage Guide

The `aws_efs_access_point` table in Steampipe provides you with information about Access Points within Amazon Elastic File System (EFS). This table enables you, as a DevOps engineer, system administrator, or other technical professional, to query access point-specific details, including the file system it is associated with, its access point ID, and other related metadata. You can utilize this table to gather insights on access points, such as their operating system type, root directory creation info, and more. The schema outlines the various attributes of the access point for you, including the access point ARN, creation time, life cycle state, and associated tags.

## Examples

### Basic info
Analyze the settings to understand the status and ownership of various access points within Amazon Elastic File System (EFS). This can help in assessing the elements within your EFS, pinpointing specific locations where changes might be needed.

```sql+postgres
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

```sql+sqlite
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
Identify the access points associated with each file system to gain insights into file ownership and root directory details. This can be useful for managing and auditing file system access within an AWS environment.

```sql+postgres
select
  name,
  access_point_id,
  file_system_id,
  owner_id,
  root_directory
from
  aws_efs_access_point
```

```sql+sqlite
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
Identify instances where access points in the AWS Elastic File System are in an error state. This could be useful in diagnosing system issues or assessing overall system health.

```sql+postgres
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

```sql+sqlite
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