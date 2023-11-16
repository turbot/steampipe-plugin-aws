---
title: "Table: aws_efs_mount_target - Query AWS EFS Mount Targets using SQL"
description: "Allows users to query AWS EFS Mount Targets for detailed information about each mount target's configuration, status, and associated resources."
---

# Table: aws_efs_mount_target - Query AWS EFS Mount Targets using SQL

The `aws_efs_mount_target` table in Steampipe provides information about mount targets within AWS Elastic File System (EFS). This table allows DevOps engineers to query mount target-specific details, including the file system ID, mount target ID, subnet ID, and security groups. Users can utilize this table to gather insights on mount targets, such as their availability, network interface, and life cycle state. The schema outlines the various attributes of the EFS mount target, including the IP address, network interface ID, owner ID, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_efs_mount_target` table, you can use the `.inspect aws_efs_mount_target` command in Steampipe.

Key columns:

- `file_system_id`: This column contains the ID of the file system for which the mount target is intended. It can be used to join this table with the `aws_efs_file_system` table.
- `mount_target_id`: This column contains the unique identifier of the mount target. It is useful for querying specific mount targets.
- `subnet_id`: This column contains the ID of the subnet in which the mount target resides. It can be used to join this table with the `aws_vpc_subnet` table.

## Examples

### Basic info

```sql
select
  mount_target_id,
  file_system_id,
  life_cycle_state,
  availability_zone_id,
  availability_zone_name
from
  aws_efs_mount_target;
```

### Get network details for each mount target

```sql
select
  mount_target_id,
  network_interface_id,
  subnet_id,
  vpc_id
from
  aws_efs_mount_target;
```
