---
title: "Steampipe Table: aws_ebs_volume_metric_read_ops_daily - Query AWS EBS Volume Metrics using SQL"
description: "Allows users to query AWS EBS Volume metrics for daily read operations."
folder: "EBS"
---

# Table: aws_ebs_volume_metric_read_ops_daily - Query AWS EBS Volume Metrics using SQL

The AWS EBS Volume Metrics is a feature of Amazon Elastic Block Store (EBS) that provides raw block-level storage that can be attached to Amazon EC2 instances. These metrics provide visibility into the performance, operation, and overall health of your volumes, allowing you to optimize usage and respond to system-wide performance changes. With the ability to query these metrics using SQL, you can gain insights into read operations on a daily basis, enhancing your ability to monitor and manage your data storage effectively.

## Table Usage Guide

The `aws_ebs_volume_metric_read_ops_daily` table in Steampipe provides you with information about the daily read operations metrics of AWS Elastic Block Store (EBS) volumes. This table allows you, as a system administrator, DevOps engineer, or other technical professional, to query details about the daily read operations performed on EBS volumes, which is useful for your performance analysis, capacity planning, and cost optimization. The schema outlines various attributes of the EBS volume metrics, including the average, maximum, and minimum read operations, as well as the sum of read operations and the time of the metric capture.

The `aws_ebs_volume_metric_read_ops_daily` table provides you with metric statistics at 24-hour intervals for the last year.

## Examples

### Basic info
Explore the performance of your AWS EBS volumes over time. This query can help you understand the volume of read operations, which can be useful in assessing system performance, identifying potential bottlenecks, and planning for capacity.

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
  aws_ebs_volume_metric_read_ops_daily
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
  aws_ebs_volume_metric_read_ops_daily
order by
  volume_id,
  timestamp;
```

### Intervals where volumes exceed 1000 average read ops
Discover the instances when the average read operations on AWS EBS volumes exceed 1000. This information can be used to identify potential performance issues or optimize resource allocation.

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
  aws_ebs_volume_metric_read_ops_daily
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
  aws_ebs_volume_metric_read_ops_daily
where
  average > 1000
order by
  volume_id,
  timestamp;
```


### Intervals where volumes exceed 8000 max read ops
Determine the instances where the daily read operations on AWS EBS volumes exceed a threshold of 8000. This can be useful in identifying potential performance issues or capacity planning for your storage infrastructure.

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
  aws_ebs_volume_metric_read_ops_daily
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
  aws_ebs_volume_metric_read_ops_daily
where
  maximum > 8000
order by
  volume_id,
  timestamp;
```


### Read, Write, and Total IOPS
Explore the average, maximum, and minimum Input/Output operations for each volume over time to understand the performance of your storage volumes. This query is useful for identifying any potential bottlenecks or inefficiencies in data transfer operations.

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