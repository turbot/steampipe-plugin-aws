---
title: "Steampipe Table: aws_ebs_volume_metric_read_ops_hourly - Query Amazon EC2 EBS Volume using SQL"
description: "Allows users to query Amazon EC2 EBS Volume Read Operations metrics on an hourly basis."
folder: "EBS"
---

# Table: aws_ebs_volume_metric_read_ops_hourly - Query Amazon EC2 EBS Volume using SQL

The AWS EBS (Elastic Block Store) Volume is a high-performance block storage service designed for use with Amazon Elastic Compute Cloud (EC2) for both throughput and transaction intensive workloads at any scale. It offers a range of volume types that are optimized to handle different workloads, including those that require high performance like transactional workloads, and those that require low cost per gigabyte like data warehousing. EBS Volumes are highly available and reliable storage volumes that can be attached to any running instance that is in the same Availability Zone.

## Table Usage Guide

The `aws_ebs_volume_metric_read_ops_hourly` table in Steampipe provides you with information about the read operations metrics of Amazon Elastic Block Store (EBS) volumes within Amazon Elastic Compute Cloud (EC2). This table allows you, as a DevOps engineer, to query volume-specific read operations details on an hourly basis, including the number of completed read operations from a volume, average, maximum, and minimum read operations, and the count of data points used for the statistical calculation. You can utilize this table to gather insights on volume performance, monitor the read activity of EBS volumes, and make data-driven decisions for performance optimization. The schema outlines the various attributes of the EBS volume read operations metrics for you, including the volume ID, timestamp, average read operations, and more.

The `aws_ebs_volume_metric_read_ops_hourly` table provides you with metric statistics at 1 hour intervals for the most recent 60 days.

## Examples

### Basic info
Explore the performance of your AWS EBS volume over time. This query allows you to track the number of read operations per hour, helping you to understand usage patterns and optimize resource allocation.

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
  aws_ebs_volume_metric_read_ops_hourly
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
  aws_ebs_volume_metric_read_ops_hourly
order by
  volume_id,
  timestamp;
```

### Intervals where volumes exceed 1000 average read ops
Identify instances where the average read operations on AWS EBS volumes exceed 1000. This can be useful in monitoring and managing resource utilization, helping to optimize performance and prevent potential bottlenecks.

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
  aws_ebs_volume_metric_read_ops_hourly
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
  aws_ebs_volume_metric_read_ops_hourly
where
  average > 1000
order by
  volume_id,
  timestamp;
```


### Intervals where volumes exceed 8000 max read ops
Identify instances where your AWS Elastic Block Store (EBS) volumes exceed 8000 maximum read operations per hour. This can help in analyzing the performance of your volumes and take necessary actions if they are under heavy load.

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
  aws_ebs_volume_metric_read_ops_hourly
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
  aws_ebs_volume_metric_read_ops_hourly
where
  maximum > 8000
order by
  volume_id,
  timestamp;
```



### Intervals where volume average iops exceeds provisioned iops
Determine the periods where the average input/output operations per second (IOPS) surpasses the provisioned IOPS for Amazon EBS volumes. This can be used to identify potential performance issues and ensure that the provisioned IOPS meets the application demand.

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
  round(r.average) + round(w.average) as iops_avg,
  round(r.average) as read_ops_avg,
  round(w.average) as write_ops_avg
from 
  aws_ebs_volume_metric_read_ops_hourly as r
join
  aws_ebs_volume_metric_write_ops_hourly as w
on 
  r.volume_id = w.volume_id
  and r.timestamp = w.timestamp
join 
  aws_ebs_volume as v
on
  v.volume_id = r.volume_id 
where 
  r.average + w.average > v.iops
order by
  r.volume_id,
  r.timestamp;
```


### Read, Write, and Total IOPS
Explore the performance of your AWS EBS volumes by evaluating the average, maximum, and minimum input/output operations per second (IOPS). This analysis can help identify any unusual activity or potential bottlenecks in your system, allowing you to optimize for better performance.

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