# Table: aws_wellarchitected_workload_share

AWS Well-Architected helps cloud architects build secure, high-performing, resilient, and efficient infrastructure for their applications and workloads. You can share a workload that you own with other AWS accounts, users, an organization, and organization units (OUs) in the same AWS Region.

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
