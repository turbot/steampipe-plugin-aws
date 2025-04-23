---
title: "Steampipe Table: aws_efs_mount_target - Query AWS EFS Mount Targets using SQL"
description: "Allows users to query AWS EFS Mount Targets for detailed information about each mount target's configuration, status, and associated resources."
folder: "EFS"
---

# Table: aws_efs_mount_target - Query AWS EFS Mount Targets using SQL

The AWS EFS Mount Target is a component of Amazon Elastic File System (EFS) that provides a network interface for a file system to connect to. It enables you to mount an Amazon EFS file system in your Amazon EC2 instance. This network interface allows the file system to connect to the network of a VPC.

## Table Usage Guide

The `aws_efs_mount_target` table in Steampipe provides you with information about mount targets within AWS Elastic File System (EFS). This table allows you, as a DevOps engineer, to query mount target-specific details, including the file system ID, mount target ID, subnet ID, and security groups. You can utilize this table to gather insights on mount targets, such as their availability, network interface, and life cycle state. The schema outlines the various attributes of the EFS mount target for you, including the IP address, network interface ID, owner ID, and associated tags.

## Examples

### Basic info
Explore the status and location of your Amazon EFS mount targets. This query is useful for understanding the availability and lifecycle state of your mount targets, which can help in optimizing resource usage and troubleshooting.

```sql+postgres
select
  mount_target_id,
  file_system_id,
  life_cycle_state,
  availability_zone_id,
  availability_zone_name
from
  aws_efs_mount_target;
```

```sql+sqlite
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
Explore the network configuration of each mount target to understand its association with different network interfaces, subnets, and virtual private clouds. This can help in assessing network-related issues and ensuring optimal configuration for enhanced performance.

```sql+postgres
select
  mount_target_id,
  network_interface_id,
  subnet_id,
  vpc_id
from
  aws_efs_mount_target;
```

```sql+sqlite
select
  mount_target_id,
  network_interface_id,
  subnet_id,
  vpc_id
from
  aws_efs_mount_target;
```