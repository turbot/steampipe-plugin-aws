---
title: "Table: aws_ebs_volume_metric_write_ops_daily - Query AWS EBS Volume Metrics using SQL"
description: "Allows users to query AWS EBS Volume Metrics for daily write operations."
---

# Table: aws_ebs_volume_metric_write_ops_daily - Query AWS EBS Volume Metrics using SQL

The `aws_ebs_volume_metric_write_ops_daily` table in Steampipe provides information about the daily write operations metrics of EBS volumes within AWS Elastic Block Store (EBS). This table allows DevOps engineers to query volume-specific details, including the number of write operations, the timestamp of data points, and the statistics for the data points. Users can utilize this table to gather insights on EBS volumes, such as the volume's write operations performance, pattern of write operations over time, and more. The schema outlines the various attributes of the EBS volume metrics, including the average, maximum, minimum, and sum of write operations, as well as the sample count for each data point.

The `aws_ebs_volume_metric_write_ops_daily` table provides metric statistics at 24 hour intervals for the last year.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ebs_volume_metric_write_ops_daily` table, you can use the `.inspect aws_ebs_volume_metric_write_ops_daily` command in Steampipe.

### Key columns:

- `title`: This column contains the title of the EBS volume. It is a key column because it can be used to join this table with other tables that contain volume-specific information.
- `average`: This column contains the average number of daily write operations for the EBS volume. It is a key column because it provides insight into the average performance of the volume.
- `timestamp`: This column contains the timestamp for each data point. It is a key column because it allows users to track the performance of the volume over time.

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
  aws_ebs_volume_metric_write_ops_daily
order by
  volume_id,
  timestamp;
```

### Intervals where volumes exceed 1000 average write ops
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
  aws_ebs_volume_metric_write_ops_daily
where
  average > 1000
order by
  volume_id,
  timestamp;
```


### Intervals where volumes exceed 8000 max write ops
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
  aws_ebs_volume_metric_write_ops_daily
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
  aws_ebs_volume_metric_read_ops_daily as r,
  aws_ebs_volume_metric_write_ops_daily as w
where 
  r.volume_id = w.volume_id
  and r.timestamp = w.timestamp
order by
  r.volume_id,
  r.timestamp;
```