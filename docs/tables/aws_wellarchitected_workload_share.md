---
title: "Table: aws_wellarchitected_workload_share - Query AWS Well-Architected Workload Share using SQL"
description: "Allows users to query AWS Well-Architected Workload Share, providing information about shared workloads within AWS Well-Architected Tool."
---

# Table: aws_wellarchitected_workload_share - Query AWS Well-Architected Workload Share using SQL

The `aws_wellarchitected_workload_share` table in Steampipe provides information about shared workloads within AWS Well-Architected Tool. This table allows DevOps engineers to query workload share-specific details, including the share ARN, workload ID, permission type, and associated metadata. Users can utilize this table to gather insights on workload shares, such as the status of the workload share, the type of permission granted, and more. The schema outlines the various attributes of the workload share, including the share ARN, workload ID, permission type, and status.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_wellarchitected_workload_share` table, you can use the `.inspect aws_wellarchitected_workload_share` command in Steampipe.

### Key columns:

- `share_arn`: The ARN of the workload share. This is the unique identifier for the workload share and can be used to join this table with others that contain workload share information.
- `workload_id`: The ID of the workload. This can be used to join this table with others that contain workload-specific information.
- `permission_type`: The permission type associated with the workload share. This can provide insights into the type of access provided through the share.

## Examples

### Basic info

```sql
select
  workload_id,
  share_id,
  shared_with
  permission_type,
  status,
  region
from
  aws_wellarchitected_workload_share;
```

### List shared workloads where invitations are pending

```sql
select
  workload_id,
  share_id,
  shared_with
  permission_type,
  status,
  region
from
  aws_wellarchitected_workload_share
where
  status = 'PENDING';
```

### List shared workloads having CONTRIBUTOR permissions

```sql
select
  workload_id,
  share_id,
  shared_with
  permission_type,
  status,
  region
from
  aws_wellarchitected_workload_share
where
  permission_type = 'CONTRIBUTOR';
```

### List shared workloads having READONLY permissions

```sql
select
  workload_id,
  share_id,
  shared_with
  permission_type,
  status,
  region
from
  aws_wellarchitected_workload_share
where
  permission_type = 'READONLY';
```
