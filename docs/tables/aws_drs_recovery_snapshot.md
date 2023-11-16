---
title: "Table: aws_drs_recovery_snapshot - Query AWS DRS Recovery Snapshot using SQL"
description: "Allows users to query DRS Recovery Snapshot data in AWS. It provides information about recovery snapshots within AWS Disaster Recovery Service (DRS). This table can be used to gather insights on recovery snapshots, including their details, associated metadata, and more."
---

# Table: aws_drs_recovery_snapshot - Query AWS DRS Recovery Snapshot using SQL

The `aws_drs_recovery_snapshot` table in Steampipe provides information about recovery snapshots within AWS Disaster Recovery Service (DRS). This table allows DevOps engineers to query snapshot-specific details, including snapshot ID, associated volume ID, start and end times, and associated metadata. Users can utilize this table to gather insights on recovery snapshots, such as snapshot status, volume size, and more. The schema outlines the various attributes of the recovery snapshot, including the snapshot ID, volume ID, start and end times, snapshot status, and volume size.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_drs_recovery_snapshot` table, you can use the `.inspect aws_drs_recovery_snapshot` command in Steampipe.

Key columns:

- `snapshot_id`: The ID of the snapshot. This is a unique identifier and can be used to join this table with other tables that contain snapshot ID information.
- `volume_id`: The ID of the volume associated with the snapshot. This can be used to join this table with other tables that contain volume ID information, providing a link between snapshots and their associated volumes.
- `start_time`: The time the snapshot was started. This can be used to analyze snapshot activity over time.

## Examples

### Basic Info

```sql
select
  snapshot_id,
  source_server_id,
  expected_timestamp,
  timestamp,
  title
from
  aws_drs_recovery_snapshot;
```

### Get source server details of each recovery snapshot

```sql
select
  r.snapshot_id,
  r.source_server_id,
  s.arn as source_server_arn,
  s.recovery_instance_id,
  s.replication_direction
from
  aws_drs_recovery_snapshot r,
  aws_drs_source_server as s
where
  r.source_server_id = s.source_server_id;
```

### Count recovery snapshots by server

```sql
select
  source_server_id,
  count(snapshot_id) as recovery_snapshot_count
from
  aws_drs_recovery_snapshot
group by
  source_server_id;
```

### List recovery snapshots taken in past 30 days

```sql
select
  snapshot_id,
  source_server_id,
  expected_timestamp,
  timestamp
from
  aws_drs_recovery_snapshot
where
  timestamp <= now() - interval '30' day;
```

### Get EBS snapshot details of a recovery snapshot

```sql
select
  r.snapshot_id,
  r.source_server_id,
  s as ebs_snapshot_id,
  e.state as snapshot_state,
  e.volume_size,
  e.volume_id,
  e.encrypted,
  e.kms_key_id,
  e.data_encryption_key_id
from
  aws_drs_recovery_snapshot as r,
  jsonb_array_elements_text(ebs_snapshots) as s,
  aws_ebs_snapshot as e
where
  r.snapshot_id = 'pit-3367d3f930778a9c3'
and
  s = e.snapshot_id;
```