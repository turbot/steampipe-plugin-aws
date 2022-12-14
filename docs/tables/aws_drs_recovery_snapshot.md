# Table: aws_drs_recovery_snapshot

An Elastic Disaster Recovery snapshot is a point-in-time copy of Amazon EBS volume, which is copied to Amazon Simple Storage Service.

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

### Get source server details for each recovery snapshot

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

<!-- Add snapshot details example query-->