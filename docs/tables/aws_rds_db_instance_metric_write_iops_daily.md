---
title: "Steampipe Table: aws_rds_db_instance_metric_write_iops_daily - Query AWS RDS DBInstance using SQL"
description: "Allows users to query AWS RDS DBInstance metrics for daily write IOPS."
folder: "RDS"
---

# Table: aws_rds_db_instance_metric_write_iops_daily - Query AWS RDS DBInstance using SQL

The AWS RDS DBInstance is a relational database service that provides you with six familiar database engines to choose from, including Amazon Aurora, PostgreSQL, MySQL, MariaDB, Oracle Database, and SQL Server. The 'write_iops_daily' metric provides the average number of disk I/O operations per second over a specified period of time. This can be used to monitor the performance of your database instance, helping to identify potential issues and optimize performance.

## Table Usage Guide

The `aws_rds_db_instance_metric_write_iops_daily` table in Steampipe provides you with information about the daily write IOPS (Input/Output Operations Per Second) metrics for each AWS RDS DBInstance. This table allows you, as a DevOps engineer, DBA, or other technical professional, to query and analyze the daily write IOPS metrics, which can be critical for your performance tuning, capacity planning, and cost management. The schema outlines the various attributes of the daily write IOPS metrics, including the DBInstance identifier, timestamp, minimum, maximum, sum, and average values, among others.

The `aws_rds_db_instance_metric_write_iops_daily` table provides you with metric statistics at 24 hour intervals for the last year.

## Examples

### Basic info
Explore the performance of your AWS RDS database instances by tracking daily write operations. This allows you to identify instances with high or low activity, helping in capacity planning and performance optimization.

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
  aws_rds_db_instance_metric_write_iops_daily
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
  aws_rds_db_instance_metric_write_iops_daily
order by
  db_instance_identifier,
  timestamp;
```

### Intervals where volumes exceed 1000 average write ops
Explore instances where the daily write operations on your AWS RDS database instances exceed an average of 1000. This helps in identifying potential performance bottlenecks and planning for capacity upgrades.

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
  aws_rds_db_instance_metric_write_iops_daily
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
  aws_rds_db_instance_metric_write_iops_daily
where
  average > 1000
order by
  db_instance_identifier,
  timestamp;
```

### Intervals where volumes exceed 8000 max write ops
Explore instances where the maximum write operations on your AWS RDS instances exceed 8000, providing insights into potential performance bottlenecks or capacity issues. This could be useful in managing resources and ensuring optimal performance of your database instances.

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
  aws_rds_db_instance_metric_write_iops_daily
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
  aws_rds_db_instance_metric_write_iops_daily
where
  maximum > 8000
order by
  db_instance_identifier,
  timestamp;
```

### Read, Write, and Total IOPS
Gain insights into the average, maximum, and minimum input/output operations per second (IOPS) for each database instance over time. This can help in understanding the performance of your databases and identifying any potential bottlenecks or areas for optimization.

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
  aws_rds_db_instance_metric_read_iops_daily as r,
  aws_rds_db_instance_metric_write_iops_daily as w
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
  aws_rds_db_instance_metric_read_iops_daily as r,
  aws_rds_db_instance_metric_write_iops_daily as w
where 
  r.db_instance_identifier = w.db_instance_identifier
  and r.timestamp = w.timestamp
order by
  r.db_instance_identifier,
  r.timestamp;
```