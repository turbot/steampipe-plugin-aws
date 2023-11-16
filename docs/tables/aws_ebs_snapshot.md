---
title: "Table: aws_ebs_snapshot - Query AWS Elastic Block Store (EBS) using SQL"
description: "Allows users to query AWS EBS snapshots, providing detailed information about each snapshot's configuration, status, and associated metadata."
---

# Table: aws_ebs_snapshot - Query AWS Elastic Block Store (EBS) using SQL

The `aws_ebs_snapshot` table in Steampipe provides information about EBS snapshots within AWS Elastic Block Store (EBS). This table allows DevOps engineers to query snapshot-specific details, including snapshot ID, description, status, volume size, and associated metadata. Users can utilize this table to gather insights on snapshots, such as snapshots with public permissions, snapshots by volume, and more. The schema outlines the various attributes of the EBS snapshot, including the snapshot ID, creation time, volume ID, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ebs_snapshot` table, you can use the `.inspect aws_ebs_snapshot` command in Steampipe.

### Key columns:

- `snapshot_id`: This is the unique identifier of the snapshot. It is useful for joining with other tables that reference snapshots by their ID.
- `volume_id`: This is the identifier of the volume from which the snapshot was created. It is useful for joining with tables that reference volumes by their ID.
- `state`: This indicates the state of the snapshot (pending, completed, or error). It is useful for filtering snapshots based on their current status.


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

### List snapshots owned by a specific AWS account

```sql
select
  snapshot_id,
  arn,
  encrypted,
  owner_id
from
  aws_ebs_snapshot
where
  owner_id = '859788737657';
```

### Get a specific snapshot by ID

```sql
select
  snapshot_id,
  arn,
  encrypted,
  owner_id
from
  aws_ebs_snapshot
where
  snapshot_id = 'snap-07bf4f91353ad71ae';
```

### List snapshots owned by Amazon (Note: This will attempt to list ALL public snapshots)

```sql
select
  snapshot_id,
  arn,
  encrypted,
  owner_id
from
  aws_ebs_snapshot
where
  owner_alias = 'amazon'
```