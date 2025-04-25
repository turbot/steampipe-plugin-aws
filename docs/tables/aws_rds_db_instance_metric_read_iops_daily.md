---
title: "Steampipe Table: aws_rds_db_instance_metric_read_iops_daily - Query AWS RDS DBInstance using SQL"
description: "Allows users to query AWS RDS DBInstance metrics for daily read IOPS (Input/Output Operations Per Second)."
folder: "RDS"
---

# Table: aws_rds_db_instance_metric_read_iops_daily - Query AWS RDS DBInstance using SQL

The AWS RDS DBInstance is a relational database service that provides you with six familiar database engines to choose from, including Amazon Aurora, PostgreSQL, MySQL, MariaDB, Oracle Database, and SQL Server. It is designed to provide a set of features to manage, scale, and operate relational databases in the cloud easily. The 'read_iops_daily' metric specifically provides insights into the average number of disk I/O operations per second for read operations in a day.

## Table Usage Guide

The `aws_rds_db_instance_metric_read_iops_daily` table in Steampipe provides you with information about the daily read IOPS metrics of AWS RDS DBInstances. This table allows you, as a DevOps engineer, to query DBInstance-specific details, including the number of read I/O operations from the DBInstance per day. You can utilize this table to gather insights on DBInstance performance, such as identifying DBInstances that have a high read I/O operations rate, which could indicate potential performance bottlenecks. The schema outlines the various attributes of the DBInstance's daily read IOPS metrics for you, including the DBInstance identifier, timestamp of the metric, and the minimum, maximum, and average number of read IOPS.

The `aws_rds_db_instance_metric_read_iops_daily` table provides you with metric statistics at 24-hour intervals for the last year.

## Examples

### Basic info
Explore the daily read input/output operations (IOPS) metrics of your Amazon RDS database instances. This can help you understand the performance trends and capacity planning for your databases over time.

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
  aws_rds_db_instance_metric_read_iops_daily
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
  aws_rds_db_instance_metric_read_iops_daily
order by
  db_instance_identifier,
  timestamp;
```

### Intervals where volumes exceed 1000 average read ops
Identify instances where the daily average read operations on AWS RDS database instances exceed 1000. This can be useful in understanding and managing resource usage and performance.

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
  aws_rds_db_instance_metric_read_iops_daily
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
  aws_rds_db_instance_metric_read_iops_daily
where
  average > 1000
order by
  db_instance_identifier,
  timestamp;
```


### Intervals where volumes exceed 8000 max read ops
This query is useful to identify periods when the read operations on your AWS RDS database instances exceed a certain threshold. It helps in monitoring and managing system performance, ensuring optimal resource utilization and preventing potential system overloads.

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
  aws_rds_db_instance_metric_read_iops_daily
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
  aws_rds_db_instance_metric_read_iops_daily
where
  maximum > 8000
order by
  db_instance_identifier,
  timestamp;
```


### Read, Write, and Total IOPS
Explore the performance of your database by analyzing the average, maximum, and minimum input/output operations per second (IOPS) for both read and write operations. This query helps in understanding the IOPS trends and can be used for capacity planning and performance tuning.

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