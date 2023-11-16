---
title: "Table: aws_ebs_volume_metric_read_ops_hourly - Query Amazon EC2 EBS Volume using SQL"
description: "Allows users to query Amazon EC2 EBS Volume Read Operations metrics on an hourly basis."
---

# Table: aws_ebs_volume_metric_read_ops_hourly - Query Amazon EC2 EBS Volume using SQL

The `aws_ebs_volume_metric_read_ops_hourly` table in Steampipe provides information about the read operations metrics of Amazon Elastic Block Store (EBS) volumes within Amazon Elastic Compute Cloud (EC2). This table allows DevOps engineers to query volume-specific read operations details on an hourly basis, including the number of completed read operations from a volume, average, maximum, and minimum read operations, and the count of data points used for the statistical calculation. Users can utilize this table to gather insights on volume performance, monitor the read activity of EBS volumes, and make data-driven decisions for performance optimization. The schema outlines the various attributes of the EBS volume read operations metrics, including the volume ID, timestamp, average read operations, and more.

The `aws_ebs_volume_metric_read_ops_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ebs_volume_metric_read_ops_hourly` table, you can use the `.inspect aws_ebs_volume_metric_read_ops_hourly` command in Steampipe.

**Key columns**:

- `volume_id`: The ID of the EBS volume. This is a key column that can be used to join this table with other tables related to EBS volumes.
- `timestamp`: The timestamp for the data point in UTC. This column is useful for tracking the read operations of a volume over time.
- `average`: The average number of completed read operations from the volume per hour. This column is important for understanding the average read performance of the volume.

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
  aws_ebs_volume_metric_read_ops_hourly
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
  aws_ebs_volume_metric_read_ops_hourly
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
  aws_ebs_volume_metric_read_ops_hourly
where
  maximum > 8000
order by
  volume_id,
  timestamp;
```



### Intervals where volume average iops exceeds provisioned iops
```sql
select 
  r.volume_id,
  r.timestamp,
  v.iops as provisioned_iops,
  round(r.average) +round(w.average) as iops_avg,
  round(r.average) as read_ops_avg,
  round(w.average) as write_ops_avg
from 
  aws_ebs_volume_metric_read_ops_hourly as r,
  aws_ebs_volume_metric_write_ops_hourly as w,
  aws_ebs_volume as v
where 
  r.volume_id = w.volume_id
  and r.timestamp = w.timestamp
  and v.volume_id = r.volume_id 
  and r.average + w.average > v.iops
order by
  r.volume_id,
  r.timestamp;
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
  aws_ebs_volume_metric_read_ops_hourly as r,
  aws_ebs_volume_metric_write_ops_hourly as w
where 
  r.volume_id = w.volume_id
  and r.timestamp = w.timestamp
order by
  r.volume_id,
  r.timestamp;
```