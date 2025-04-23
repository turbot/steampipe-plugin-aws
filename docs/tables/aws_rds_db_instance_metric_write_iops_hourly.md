---
title: "Steampipe Table: aws_rds_db_instance_metric_write_iops_hourly - Query AWS RDS DBInstance Metrics using SQL"
description: "Allows users to query AWS RDS DBInstance write IOPS metrics on an hourly basis."
folder: "RDS"
---

# Table: aws_rds_db_instance_metric_write_iops_hourly - Query AWS RDS DBInstance Metrics using SQL

The AWS RDS DBInstance Metrics is a feature of Amazon RDS that provides metrics data for a DB instance. It allows you to monitor and manage the performance of the DB instance by providing data in a readable, user-friendly format. It includes metrics such as Write IOPS, which represents the average number of disk I/O operations per second.

## Table Usage Guide

The `aws_rds_db_instance_metric_write_iops_hourly` table in Steampipe provides you with information about the Input/Output operations per second (IOPS) for write operations on an AWS RDS DBInstance, aggregated on an hourly basis. You can use this table to query DBInstance-specific details, including the number of write IOPS, the timestamp of the data point, and the statistical value. This table allows you, as a DevOps engineer, database administrator, or other technical professional, to gather insights on the write performance of your DBInstances. You can identify periods of high write activity, monitor the impact of performance tuning measures, and more. The schema outlines the various attributes of the DBInstance metric for you, including the DBInstance identifier, the period, the unit, and the timestamp.

The `aws_rds_db_instance_metric_write_iops_hourly` table provides you with metric statistics at 1 hour intervals for the most recent 60 days.

## Examples

### Basic info
Explore the performance of your AWS RDS instances by tracking hourly write operations. This allows for proactive management and optimization of database performance.

```sql+postgres
select
  db_instance_identifier,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_rds_db_instance_metric_write_iops_hourly
order by
  db_instance_identifier,
  timestamp;
```

```sql+sqlite
select
  db_instance_identifier,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_rds_db_instance_metric_write_iops_hourly
order by
  db_instance_identifier,
  timestamp;
```

### Intervals where volumes exceed 1000 average write ops
Identify instances where the average write operations per hour exceed 1000 in your AWS RDS database instances. This can help in detecting high usage periods and planning for capacity upgrades.

```sql+postgres
select
  db_instance_identifier,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_rds_db_instance_metric_write_iops_hourly
where
  average > 1000
order by
  db_instance_identifier,
  timestamp;
```

```sql+sqlite
select
  db_instance_identifier,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_rds_db_instance_metric_write_iops_hourly
where
  average > 1000
order by
  db_instance_identifier,
  timestamp;
```

### Intervals where volumes exceed 8000 max write ops
Identify instances where database write operations exceed a specified threshold. This is useful for monitoring system performance and identifying potential bottlenecks or periods of heavy load.

```sql+postgres
select
  db_instance_identifier,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_rds_db_instance_metric_write_iops_hourly
where
  maximum > 8000
order by
  db_instance_identifier,
  timestamp;
```

```sql+sqlite
select
  db_instance_identifier,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_rds_db_instance_metric_write_iops_hourly
where
  maximum > 8000
order by
  db_instance_identifier,
  timestamp;
```

### Intervals where volume average iops exceeds provisioned iops
Identify instances where the average input/output operations per second (IOPS) exceeds the provisioned IOPS. This helps in monitoring the performance and ensuring the efficient use of resources in your database environment.

```sql+postgres
select 
  r.db_instance_identifier,
  r.timestamp,
  v.iops as provisioned_iops,
  round(r.average) +round(w.average) as iops_avg,
  round(r.average) as read_ops_avg,
  round(w.average) as write_ops_avg
from 
  aws_rds_db_instance_metric_read_iops_hourly as r,
  aws_rds_db_instance_metric_write_iops_hourly as w,
  aws_rds_db_instance as v
where 
  r.db_instance_identifier = w.db_instance_identifier
  and r.timestamp = w.timestamp
  and v.db_instance_identifier = r.db_instance_identifier 
  and r.average + w.average > v.iops
order by
  r.db_instance_identifier,
  r.timestamp;
```

```sql+sqlite
select 
  r.db_instance_identifier,
  r.timestamp,
  v.iops as provisioned_iops,
  round(r.average) + round(w.average) as iops_avg,
  round(r.average) as read_ops_avg,
  round(w.average) as write_ops_avg
from 
  aws_rds_db_instance_metric_read_iops_hourly as r
join
  aws_rds_db_instance_metric_write_iops_hourly as w
on
  r.db_instance_identifier = w.db_instance_identifier
  and r.timestamp = w.timestamp
join
  aws_rds_db_instance as v
on
  v.db_instance_identifier = r.db_instance_identifier 
where 
  r.average + w.average > v.iops
order by
  r.db_instance_identifier,
  r.timestamp;
```


### Read, Write, and Total IOPS
This query enables you to monitor the performance of your AWS RDS instances by providing insights into the average, maximum, and minimum Input/Output operations per second (IOPS). By analyzing these metrics, you can optimize your database performance, identify potential bottlenecks, and make informed decisions about capacity planning.

```sql+postgres
select 
  r.db_instance_identifier,
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
  aws_rds_db_instance_metric_read_iops_hourly as r,
  aws_rds_db_instance_metric_write_iops_hourly as w
where 
  r.db_instance_identifier = w.db_instance_identifier
  and r.timestamp = w.timestamp
order by
  r.db_instance_identifier,
  r.timestamp;
```

```sql+sqlite
select 
  r.db_instance_identifier,
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
  aws_rds_db_instance_metric_read_iops_hourly as r,
  aws_rds_db_instance_metric_write_iops_hourly as w
where 
  r.db_instance_identifier = w.db_instance_identifier
  and r.timestamp = w.timestamp
order by
  r.db_instance_identifier,
  r.timestamp;
```