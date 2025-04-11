---
title: "Steampipe Table: aws_ebs_volume_metric_write_ops_hourly - Query AWS EBS Volume Metrics using SQL"
description: "Allows users to query AWS EBS Volume Metrics on hourly write operations."
folder: "EBS"
---

# Table: aws_ebs_volume_metric_write_ops_hourly - Query AWS EBS Volume Metrics using SQL

The AWS EBS (Elastic Block Store) Volume Metrics is a feature that allows you to monitor the performance of your EBS volumes for analysis and troubleshooting. With the 'write_ops' metric, you can track the number of write operations performed on a specified EBS volume per hour. This data can be queried using SQL, providing an accessible way to monitor and manage the performance of your EBS volumes.

## Table Usage Guide

The `aws_ebs_volume_metric_write_ops_hourly` table in Steampipe provides you with information about the hourly write operations metrics of AWS Elastic Block Store (EBS) volumes. This table allows you, as a cloud engineer, a member of a DevOps team, or a data analyst, to query and analyze the hourly write operation details of EBS volumes, including the number of write operations and the timestamp of the data points. You can utilize this table to track write operations, monitor EBS performance, and plan capacity. The schema outlines the various attributes of the EBS volume metrics for you, including the volume ID, timestamp, and the number of write operations.

The `aws_ebs_volume_metric_write_ops_hourly` table provides you with metric statistics at 1 hour intervals for the most recent 60 days.

## Examples

### Basic info
Gain insights into the performance of your AWS EBS volumes by analyzing write operations over time. This can assist in identifying potential issues, optimizing resource usage, and planning capacity.

```sql+postgres
select
  volume_id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_ebs_volume_metric_write_ops_hourly
order by
  volume_id,
  timestamp;
```

```sql+sqlite
select
  volume_id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_ebs_volume_metric_write_ops_hourly
order by
  volume_id,
  timestamp;
```

### Intervals where volumes exceed 1000 average write ops
Discover the instances where the average write operations on your AWS EBS volumes exceed 1000 per hour. This can be useful to identify potential performance issues or unusual activity on your volumes.

```sql+postgres
select
  volume_id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_ebs_volume_metric_write_ops_hourly
where
  average > 1000
order by
  volume_id,
  timestamp;
```

```sql+sqlite
select
  volume_id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_ebs_volume_metric_write_ops_hourly
where
  average > 1000
order by
  volume_id,
  timestamp;
```


### Intervals where volumes exceed 8000 max write ops
Identify instances where the maximum write operations on AWS EBS volumes exceed 8000 within an hour. This can help monitor and manage storage performance, ensuring optimal operation and preventing potential issues.

```sql+postgres
select
  volume_id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_ebs_volume_metric_write_ops_hourly
where
  maximum > 8000
order by
  volume_id,
  timestamp;
```

```sql+sqlite
select
  volume_id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_ebs_volume_metric_write_ops_hourly
where
  maximum > 8000
order by
  volume_id,
  timestamp;
```

### Intervals where volume average iops exceeds provisioned iops
Identify instances where the average input/output operations per second (IOPS) surpasses the provisioned IOPS on your AWS EBS volumes. This is crucial for optimizing your storage performance and preventing any potential bottlenecks.

```sql+postgres
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

```sql+sqlite
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
Analyze the settings to understand the average, maximum, and minimum input/output operations per second (IOPS) for both read and write operations on AWS EBS volumes. This helps in assessing the performance and identifying any potential bottlenecks in data transfer.

```sql+postgres
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

```sql+sqlite
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