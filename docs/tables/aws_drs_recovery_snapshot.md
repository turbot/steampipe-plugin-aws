---
title: "Steampipe Table: aws_drs_recovery_snapshot - Query AWS DRS Recovery Snapshot using SQL"
description: "Allows users to query DRS Recovery Snapshot data in AWS. It provides information about recovery snapshots within AWS Disaster Recovery Service (DRS). This table can be used to gather insights on recovery snapshots, including their details, associated metadata, and more."
folder: "Elastic Disaster Recovery (DRS)"
---

# Table: aws_drs_recovery_snapshot - Query AWS DRS Recovery Snapshot using SQL

The AWS Disaster Recovery Service (DRS) Recovery Snapshot is a feature of AWS DRS that allows you to capture the state of your resources at a specific point in time. This is crucial for disaster recovery purposes, enabling you to restore your system to a previous state in case of a disaster. The snapshot includes all of your data, applications, and configurations, providing a comprehensive backup of your resources.

## Table Usage Guide

The `aws_drs_recovery_snapshot` table in Steampipe provides you with information about recovery snapshots within AWS Disaster Recovery Service (DRS). This table enables you, as a DevOps engineer, to query snapshot-specific details, including snapshot ID, associated volume ID, start and end times, and associated metadata. You can utilize this table to gather insights on recovery snapshots, such as snapshot status, volume size, and more. The schema outlines the various attributes of the recovery snapshot for you, including the snapshot ID, volume ID, start and end times, snapshot status, and volume size.

## Examples

### Basic Info
Discover the segments that require recovery snapshots in your AWS Disaster Recovery Service. This can help you anticipate and manage potential system recovery needs effectively.

```sql+postgres
select
  snapshot_id,
  source_server_id,
  expected_timestamp,
  timestamp,
  title
from
  aws_drs_recovery_snapshot;
```

```sql+sqlite
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
This query is useful for gaining insights into the origin of each recovery snapshot in a disaster recovery system. It allows users to identify the specific source server of each snapshot, which can be beneficial for system audits, recovery planning, and troubleshooting.

```sql+postgres
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

```sql+sqlite
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
Determine the quantity of recovery snapshots for each server to understand the frequency of data recovery measures taken. This is useful for assessing the robustness of your data backup strategy.

```sql+postgres
select
  source_server_id,
  count(snapshot_id) as recovery_snapshot_count
from
  aws_drs_recovery_snapshot
group by
  source_server_id;
```

```sql+sqlite
select
  source_server_id,
  count(snapshot_id) as recovery_snapshot_count
from
  aws_drs_recovery_snapshot
group by
  source_server_id;
```

### List recovery snapshots taken in past 30 days
Identify instances where recovery snapshots have been taken in the past 30 days. This is useful for maintaining an up-to-date backup and recovery strategy in your AWS environment.

```sql+postgres
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

```sql+sqlite
select
  snapshot_id,
  source_server_id,
  expected_timestamp,
  timestamp
from
  aws_drs_recovery_snapshot
where
  timestamp <= datetime('now', '-30 day');
```

### Get EBS snapshot details of a recovery snapshot
Determine the specifics of a particular recovery snapshot within your AWS Disaster Recovery service, such as its state, volume size, and encryption details. This can be useful for understanding the properties of your recovery snapshots and ensuring they meet your data security and storage requirements.

```sql+postgres
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

```sql+sqlite
select
  r.snapshot_id,
  r.source_server_id,
  json_extract(s.value, '$') as ebs_snapshot_id,
  e.state as snapshot_state,
  e.volume_size,
  e.volume_id,
  e.encrypted,
  e.kms_key_id,
  e.data_encryption_key_id
from
  aws_drs_recovery_snapshot as r,
  json_each(r.ebs_snapshots) as s,
  aws_ebs_snapshot as e
where
  r.snapshot_id = 'pit-3367d3f930778a9c3'
and
  json_extract(s.value, '$') = e.snapshot_id;
```