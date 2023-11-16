---
title: "Table: aws_ebs_volume_metric_read_ops - Query AWS EBS Volume using SQL"
description: "Allows users to query AWS EBS Volume read operations metrics."
---

# Table: aws_ebs_volume_metric_read_ops - Query AWS EBS Volume using SQL

The `aws_ebs_volume_metric_read_ops` table in Steampipe provides information about read operations metrics of volumes within AWS Elastic Block Store (EBS). This table allows DevOps engineers to query volume-specific details, including the number of read operations that completed, the timestamp of the measurement, and associated metadata. Users can utilize this table to gather insights on volumes, such as the frequency of read operations, the performance of volumes over time, and more. The schema outlines the various attributes of the EBS volume read operations metrics, including the volume id, timestamp, and the number of read operations.

The `aws_ebs_volume_metric_read_ops` table provides metric statistics at 5 minute intervals for the most recent 5 days.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ebs_volume_metric_read_ops` table, you can use the `.inspect aws_ebs_volume_metric_read_ops` command in Steampipe.

Key columns:

- `title`: This is the title of the EBS volume metric. It is important as it identifies the specific volume that the read operation metrics pertain to.
- `value`: This column represents the number of read operations that completed. It is useful for monitoring the read operations of the volume over time.
- `timestamp`: This column records the time when the measurement was taken. It is crucial for tracking the performance of the volume over time.

## Examples

### Basic info

```sql
select
  volume_id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_ebs_volume_metric_read_ops
order by
  volume_id,
  timestamp;
```

### Intervals where volumes exceed 1000 average read ops
```sql
select
  volume_id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_ebs_volume_metric_read_ops
where
  average > 1000
order by
  volume_id,
  timestamp;
```


### Intervals where volumes exceed 8000 max read ops
```sql
select
  volume_id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_ebs_volume_metric_read_ops
where
  maximum > 8000
order by
  volume_id,
  timestamp;
```


### Read, Write, and Total IOPS

```sql
select 
  r.volume_id,
  r.timestamp,
  round(r.average) + round(w.average) as iops_avg,
  round(r.average) as read_ops_avg,
  round(w.average) as write_ops_avg,
  round(r.maximum) + round(w.maximum) as iops_max,
  round(r.maximum) as read_ops_max,
  round(w.maximum) as write_ops_max,
  round(r.minimum) + round(w.minimum) as iops_min,
  round(r.minimum) as read_ops_min,
  round(w.minimum) as write_ops_min
from 
  aws_ebs_volume_metric_read_ops as r,
  aws_ebs_volume_metric_write_ops as w
where 
  r.volume_id = w.volume_id
  and r.timestamp = w.timestamp
order by
  r.volume_id,
  r.timestamp;
```