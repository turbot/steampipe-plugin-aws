---
title: "Table: aws_ebs_volume_metric_write_ops - Query AWS Elastic Block Store (EBS) using SQL"
description: "Allows users to query AWS Elastic Block Store (EBS) volume write operations metrics."
---

# Table: aws_ebs_volume_metric_write_ops - Query AWS Elastic Block Store (EBS) using SQL

The `aws_ebs_volume_metric_write_ops` table in Steampipe provides information about the write operations metrics of EBS volumes within AWS Elastic Block Store (EBS). This table allows DevOps engineers to query volume-specific details, including the number of write operations, the timestamp of the data point, and the statistical value of the data point. Users can utilize this table to gather insights on EBS volumes, such as volume performance, write load, and more. The schema outlines the various attributes of the EBS volume write operations metrics, including the volume ID, timestamp, and statistical values.

The `aws_ebs_volume_metric_write_ops` table provides metric statistics at 5 minute intervals for the most recent 5 days.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ebs_volume_metric_write_ops` table, you can use the `.inspect aws_ebs_volume_metric_write_ops` command in Steampipe.

Key columns:

- `title_id`: The ID of the EBS volume. This can be used to join this table with other tables that contain EBS volume information.
- `timestamp`: The timestamp for the data point. This can be used to track the write operations over time.
- `average`: The statistical value for the data point. This can provide insights into the average number of write operations for the EBS volume.

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
  aws_ebs_volume_metric_write_ops
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
  aws_ebs_volume_metric_write_ops
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
  aws_ebs_volume_metric_write_ops
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