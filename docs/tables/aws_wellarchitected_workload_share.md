---
title: "Steampipe Table: aws_wellarchitected_workload_share - Query AWS Well-Architected Workload Share using SQL"
description: "Allows users to query AWS Well-Architected Workload Share, providing information about shared workloads within AWS Well-Architected Tool."
folder: "Well-Architected"
---

# Table: aws_wellarchitected_workload_share - Query AWS Well-Architected Workload Share using SQL

The AWS Well-Architected Workload Share is a feature of the AWS Well-Architected Tool that allows you to share your workloads with other AWS accounts. This function enables collaboration with others to review and improve the design and architecture of your applications. The sharing process is secure and managed, ensuring only authorized access to your workload information.

## Table Usage Guide

The `aws_wellarchitected_workload_share` table in Steampipe provides you with information about shared workloads within AWS Well-Architected Tool. This table allows you, as a DevOps engineer, to query workload share-specific details, including the share ARN, workload ID, permission type, and associated metadata. You can utilize this table to gather insights on workload shares, such as the status of the workload share, the type of permission granted, and more. The schema outlines the various attributes of the workload share for you, including the share ARN, workload ID, permission type, and status.

## Examples

### Basic info
Explore which workload shares in your AWS Well-Architected Tool have been shared with others, their permission types, and status. This can help you manage and control access to your workloads effectively across different regions.

```sql+postgres
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

```sql+sqlite
select
  workload_id,
  share_id,
  shared_with,
  permission_type,
  status,
  region
from
  aws_wellarchitected_workload_share;
```

### List shared workloads where invitations are pending
Determine the areas in which workload shares in AWS Well-Architected Tool are still pending approval. This can be useful for managing workload collaborations and ensuring timely access for all involved parties.

```sql+postgres
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

```sql+sqlite
select
  workload_id,
  share_id,
  shared_with,
  permission_type,
  status,
  region
from
  aws_wellarchitected_workload_share
where
  status = 'PENDING';
```

### List shared workloads having CONTRIBUTOR permissions
Identify shared workloads where the user has been granted 'Contributor' permissions. This can be useful in managing access rights and understanding the distribution of workload responsibilities.

```sql+postgres
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

```sql+sqlite
select
  workload_id,
  share_id,
  shared_with,
  permission_type,
  status,
  region
from
  aws_wellarchitected_workload_share
where
  permission_type = 'CONTRIBUTOR';
```

### List shared workloads having READONLY permissions
Identify shared workloads that have been granted 'READONLY' permissions. This allows you to understand which external entities have limited access to your workloads, helping to maintain security and control over your AWS environment.

```sql+postgres
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

```sql+sqlite
select
  workload_id,
  share_id,
  shared_with,
  permission_type,
  status,
  region
from
  aws_wellarchitected_workload_share
where
  permission_type = 'READONLY';
```