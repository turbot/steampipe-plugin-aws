---
title: "Steampipe Table: aws_ebs_volume_metric_write_ops_daily - Query AWS EBS Volume Metrics using SQL"
description: "Allows users to query AWS EBS Volume Metrics for daily write operations."
folder: "EBS"
---

# Table: aws_ebs_volume_metric_write_ops_daily - Query AWS EBS Volume Metrics using SQL

The AWS EBS Volume Metrics provides a way to monitor the performance of your Amazon Elastic Block Store (EBS) volumes. It allows you to capture write operations on a daily basis, which can help you optimize your storage usage. These metrics can be queried using SQL, providing a flexible and efficient way to analyze your EBS performance data.

## Table Usage Guide

The `aws_ebs_volume_metric_write_ops_daily` table in Steampipe provides you with information about the daily write operations metrics of EBS volumes within AWS Elastic Block Store (EBS). This table allows you, as a DevOps engineer, to query volume-specific details, including the number of write operations, the timestamp of data points, and the statistics for the data points. You can utilize this table to gather insights on EBS volumes, such as the volume's write operations performance, pattern of write operations over time, and more. The schema outlines the various attributes of the EBS volume metrics for you, including the average, maximum, minimum, and sum of write operations, as well as the sample count for each data point.

The `aws_ebs_volume_metric_write_ops_daily` table provides you with metric statistics at 24 hour intervals for the last year.

## Examples

### Basic info
This query allows you to analyze the daily write operations of AWS EBS volumes. It can be used to gain insights into the performance and usage patterns of your volumes, helping optimize resource allocation and troubleshoot potential issues.

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
  aws_ebs_volume_metric_write_ops_daily
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
  aws_ebs_volume_metric_write_ops_daily
order by
  volume_id,
  timestamp;
```

### Intervals where volumes exceed 1000 average write ops
Identify instances where the daily average write operations on AWS EBS volumes exceed 1000. This is useful for monitoring usage patterns and potentially preventing system overloads.
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
  aws_ebs_volume_metric_write_ops_daily
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
  aws_ebs_volume_metric_write_ops_daily
where
  average > 1000
order by
  volume_id,
  timestamp;
```


### Intervals where volumes exceed 8000 max write ops
Determine the instances where the maximum write operations on AWS EBS volumes surpass the 8000 mark. This is useful to identify potential bottlenecks in your storage system and take proactive measures to prevent performance degradation.
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
  aws_ebs_volume_metric_write_ops_daily
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
  aws_ebs_volume_metric_write_ops_daily
where
  maximum > 8000
order by
  volume_id,
  timestamp;
```


### Read, Write, and Total IOPS
Explore the performance of your AWS EBS volumes by understanding their input/output operations over time. This query will help you analyze the average, maximum, and minimum read/write operations, allowing you to optimize your storage usage and troubleshoot any potential issues.

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
  aws_ebs_volume_metric_read_ops_daily as r,
  aws_ebs_volume_metric_write_ops_daily as w
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
  aws_ebs_volume_metric_read_ops_daily as r,
  aws_ebs_volume_metric_write_ops_daily as w
where 
  r.volume_id = w.volume_id
  and r.timestamp = w.timestamp
order by
  r.volume_id,
  r.timestamp;
```