---
title: "Table: aws_ebs_volume_metric_read_ops_daily - Query AWS EBS Volume Metrics using SQL"
description: "Allows users to query AWS EBS Volume metrics for daily read operations."
---

# Table: aws_ebs_volume_metric_read_ops_daily - Query AWS EBS Volume Metrics using SQL

The `aws_ebs_volume_metric_read_ops_daily` table in Steampipe provides information about the daily read operations metrics of AWS Elastic Block Store (EBS) volumes. This table allows system administrators, DevOps engineers, and other technical professionals to query details about the daily read operations performed on EBS volumes, which is useful for performance analysis, capacity planning, and cost optimization. The schema outlines various attributes of the EBS volume metrics, including the average, maximum, and minimum read operations, as well as the sum of read operations and the time of the metric capture.

The `aws_ebs_volume_metric_read_ops_daily` table provides metric statistics at 24 hour intervals for the last year.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ebs_volume_metric_read_ops_daily` table, you can use the `.inspect aws_ebs_volume_metric_read_ops_daily` command in Steampipe.

### Key columns:

- `title_id`: The ID of the EBS volume. It is a key column for joining with other tables to get more information about the EBS volume.
- `timestamp`: The timestamp for the metric data point. It is essential for tracking the volume's performance over time.
- `average`: The average number of read operations during the specified period. It provides a baseline for comparing volume performance.

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
  aws_ebs_volume_metric_read_ops_daily
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
  aws_ebs_volume_metric_read_ops_daily
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
  aws_ebs_volume_metric_read_ops_daily
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