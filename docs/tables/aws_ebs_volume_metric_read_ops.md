---
title: "Steampipe Table: aws_ebs_volume_metric_read_ops - Query AWS EBS Volume using SQL"
description: "Allows users to query AWS EBS Volume read operations metrics."
folder: "EBS"
---

# Table: aws_ebs_volume_metric_read_ops - Query AWS EBS Volume using SQL

The AWS EBS Volume is a block-level storage device that you can attach to a single EC2 instance. It allows you to persist data past the lifespan of a single Amazon EC2 instance, and the data on an EBS volume is replicated within its availability zone to prevent data loss due to failure. The 'read_ops' metric provides the total number of read operations from an EBS volume, which can be queried using SQL.

## Table Usage Guide

The `aws_ebs_volume_metric_read_ops` table in Steampipe provides you with information about read operations metrics of volumes within AWS Elastic Block Store (EBS). This table allows you, as a DevOps engineer, to query volume-specific details, including the number of read operations that have been completed, the timestamp of the measurement, and associated metadata. You can utilize this table to gather insights on volumes, such as the frequency of read operations, the performance of volumes over time, and more. The schema outlines the various attributes of the EBS volume read operations metrics for you, including the volume id, timestamp, and the number of read operations.

The `aws_ebs_volume_metric_read_ops` table provides you with metric statistics at 5-minute intervals for the most recent 5 days.

## Examples

### Basic info
Analyze the performance of your Amazon EBS volumes over time. This query aids in understanding the read operation metrics, including minimum, maximum, average and total read operations, helping you optimize your resource usage and troubleshoot potential issues.

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
  aws_ebs_volume_metric_read_ops
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
  aws_ebs_volume_metric_read_ops
order by
  volume_id,
  timestamp;
```

### Intervals where volumes exceed 1000 average read ops
Determine the periods when the average read operations on AWS EBS volumes surpass a certain threshold. This is useful for identifying potential performance bottlenecks and planning for capacity upgrades.

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
  aws_ebs_volume_metric_read_ops
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
  aws_ebs_volume_metric_read_ops
where
  average > 1000
order by
  volume_id,
  timestamp;
```


### Intervals where volumes exceed 8000 max read ops
Determine the intervals where the operation count on your Elastic Block Store (EBS) volumes exceeds 8000 read operations. This can assist in identifying potential performance issues or bottlenecks in your AWS environment.

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
  aws_ebs_volume_metric_read_ops
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
  aws_ebs_volume_metric_read_ops
where
  maximum > 8000
order by
  volume_id,
  timestamp;
```


### Read, Write, and Total IOPS
Assess the performance of your AWS EBS volumes by examining the average, maximum, and minimum Input/Output Operations Per Second (IOPS). This can help you understand your application’s load on the volumes and plan for capacity or performance improvements.

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
  aws_ebs_volume_metric_read_ops as r,
  aws_ebs_volume_metric_write_ops as w
where 
  r.volume_id = w.volume_id
  and r.timestamp = w.timestamp
order by
  r.volume_id,
  r.timestamp;
```