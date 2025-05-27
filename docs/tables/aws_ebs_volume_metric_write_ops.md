---
title: "Steampipe Table: aws_ebs_volume_metric_write_ops - Query AWS Elastic Block Store (EBS) using SQL"
description: "Allows users to query AWS Elastic Block Store (EBS) volume write operations metrics."
folder: "EBS"
---

# Table: aws_ebs_volume_metric_write_ops - Query AWS Elastic Block Store (EBS) using SQL

The AWS Elastic Block Store (EBS) is a high-performance block storage service designed for use with Amazon EC2 for both throughput and transaction intensive workloads at any scale. It provides persistent block level storage volumes for use with EC2 instances. The "write_ops" metric represents the number of write operations performed on the EBS volume.

## Table Usage Guide

The `aws_ebs_volume_metric_write_ops` table in Steampipe provides you with information about the write operations metrics of EBS volumes within AWS Elastic Block Store (EBS). This table allows you, as a DevOps engineer, to query volume-specific details, including the number of write operations, the timestamp of the data point, and the statistical value of the data point. You can utilize this table to gather insights on EBS volumes, such as volume performance, write load, and more. The schema outlines the various attributes of the EBS volume write operations metrics for you, including the volume ID, timestamp, and statistical values.

The `aws_ebs_volume_metric_write_ops` table provides you with metric statistics at 5 minute intervals for the most recent 5 days.

## Examples

### Basic info
Gain insights into the performance of your AWS EBS volumes over time. This query helps in monitoring the write operations, which aids in identifying potential bottlenecks or performance issues.

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
  aws_ebs_volume_metric_write_ops
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
  aws_ebs_volume_metric_write_ops
order by
  volume_id,
  timestamp;
```

### Intervals where volumes exceed 1000 average write ops
Identify instances where the average write operations on AWS EBS volumes exceed 1000. This can be useful in monitoring performance and identifying potential bottlenecks or areas for optimization.

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
  aws_ebs_volume_metric_write_ops
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
  aws_ebs_volume_metric_write_ops
where
  average > 1000
order by
  volume_id,
  timestamp;
```


### Intervals where volumes exceed 8000 max write ops
Identify instances where the maximum write operations on AWS EBS volumes exceed 8000. This can be useful in understanding the load on your EBS volumes, and may help you optimize your resources for better performance.

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
  aws_ebs_volume_metric_write_ops
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
  aws_ebs_volume_metric_write_ops
where
  maximum > 8000
order by
  volume_id,
  timestamp;
```


### Read, Write, and Total IOPS
Explore the performance of your storage volumes by analyzing the average, maximum, and minimum Input/Output operations per second (IOPS). This allows you to monitor and optimize your storage efficiency, ensuring smooth operations.

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
  aws_ebs_volume_metric_read_ops as r,
  aws_ebs_volume_metric_write_ops as w
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
  aws_ebs_volume_metric_read_ops as r
join
  aws_ebs_volume_metric_write_ops as w on r.volume_id = w.volume_id and r.timestamp = w.timestamp
order by
  r.volume_id,
  r.timestamp;
```