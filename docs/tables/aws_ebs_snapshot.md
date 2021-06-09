# Table: aws_ebs_snapshot

An EBS snapshot is a point-in-time copy of Amazon EBS volume, which is copied to Amazon Simple Storage Service.

## Examples

### List of snapshots which are not encrypted

```sql
select
  snapshot_id,
  arn,
  encrypted
from
  aws_ebs_snapshot
where
  not encrypted;
```

### List of EBS snapshots which are publicly accessible

```sql
select
  snapshot_id,
  arn,
  volume_id,
  perm ->> 'UserId' as userid,
  perm ->> 'Group' as group
from
  aws_ebs_snapshot
  cross join jsonb_array_elements(create_volume_permissions) as perm
where
  perm ->> 'Group' = 'all';
```

### Find the Account IDs with which the snapshots are shared

```sql
select
  snapshot_id,
  volume_id,
  perm ->> 'UserId' as account_ids
from
  aws_ebs_snapshot
  cross join jsonb_array_elements(create_volume_permissions) as perm;
```

### Find the snapshot count per volume

```sql
select
  volume_id,
  count(snapshot_id) as snapshot_id
from
  aws_ebs_snapshot
group by
  volume_id;
```
