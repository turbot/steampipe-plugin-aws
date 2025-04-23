---
title: "Steampipe Table: aws_rds_db_instance_metric_write_iops - Query AWS RDS DBInstance Write IOPS using SQL"
description: "Allows users to query AWS RDS DBInstance Write IOPS to retrieve metrics on the write input/output operations per second."
folder: "RDS"
---

# Table: aws_rds_db_instance_metric_write_iops - Query AWS RDS DBInstance Write IOPS using SQL

The AWS RDS DBInstance Write IOPS is a metric that measures the average number of disk I/O operations per second for write operations. It is part of Amazon RDS (Relational Database Service) which provides resizable capacity for an industry-standard relational database and manages common database administration tasks. Monitoring this metric can help in understanding the performance of the database and identifying any potential issues related to data write operations.

## Table Usage Guide

The `aws_rds_db_instance_metric_write_iops` table in Steampipe gives you information about the Write IOPS (Input/Output Operations Per Second) metrics of AWS RDS DBInstances. This table allows you, as a DevOps engineer, to query details related to the write operations performance of your RDS DBInstances, including the average, maximum, and minimum values for a specified period. You can utilize this table to monitor and analyze the write performance of your DBInstances, helping you optimize the performance and reliability of your database operations. The schema outlines the various attributes of the Write IOPS metrics, including the DBInstance identifier, timestamp, and the statistics for the period.

The `aws_rds_db_instance_metric_write_iops` table provides metric statistics at 5 minute intervals for the most recent 5 days.

## Examples

### Basic info
Explore the performance of various AWS RDS database instances over time. This query is useful for identifying trends or anomalies in write operations, which can inform optimization strategies or troubleshooting efforts.

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
  aws_rds_db_instance_metric_write_iops
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
  aws_rds_db_instance_metric_write_iops
order by
  db_instance_identifier,
  timestamp;
```

### Intervals where volumes exceed 1000 average write ops
Identify instances where the average write operations on your AWS RDS database instances exceed 1000. This can help you pinpoint periods of high database activity and assess the need for potential performance optimization or resource scaling.

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
  aws_rds_db_instance_metric_write_iops
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
  aws_rds_db_instance_metric_write_iops
where
  average > 1000
order by
  db_instance_identifier,
  timestamp;
```

### Intervals where volumes exceed 8000 max write ops
Identify instances where the maximum write operations on your AWS RDS database instances exceed 8000. This can help in assessing the load on your databases and planning for capacity upgrades or performance optimization.

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
  aws_rds_db_instance_metric_write_iops
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
  aws_rds_db_instance_metric_write_iops
where
  maximum > 8000
order by
  db_instance_identifier,
  timestamp;
```

### Read, Write, and Total IOPS
Analyze the settings to understand the average, maximum, and minimum input/output operations per second (IOPS) for each database instance. This can help pinpoint performance bottlenecks and optimize database operations for better efficiency.

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
  aws_rds_db_instance_metric_read_iops as r,
  aws_rds_db_instance_metric_write_iops as w
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
  aws_rds_db_instance_metric_read_iops as r,
  aws_rds_db_instance_metric_write_iops as w
where 
  r.db_instance_identifier = w.db_instance_identifier
  and r.timestamp = w.timestamp
order by
  r.db_instance_identifier,
  r.timestamp;
```