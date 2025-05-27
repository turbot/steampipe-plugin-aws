---
title: "Steampipe Table: aws_rds_db_instance_metric_read_iops - Query AWS RDS DBInstanceMetricReadIops using SQL"
description: "Allows users to query AWS RDS DBInstanceMetricReadIops to retrieve and monitor the read IOPS (Input/Output Operations Per Second) metrics for Amazon RDS DB instances."
folder: "RDS"
---

# Table: aws_rds_db_instance_metric_read_iops - Query AWS RDS DBInstanceMetricReadIops using SQL

The AWS RDS DB Instance Metric Read IOPS is a performance metric for Amazon Relational Database Service (RDS) that measures the average number of disk I/O operations per second for read operations. This metric is useful to monitor the read activity on your RDS DB instance and can help you identify potential performance issues. It is part of the suite of CloudWatch metrics for RDS that provides detailed visibility into the health, performance, and availability of your RDS databases.

## Table Usage Guide

The `aws_rds_db_instance_metric_read_iops` table in Steampipe provides you with information about the read IOPS metrics of AWS RDS DB instances. This table allows you, as a DevOps engineer or database administrator, to query and monitor the read IOPS metrics, which can be useful for performance tuning and capacity planning. The read IOPS refers to the number of read input/output operations per second. The schema outlines the various attributes of the DB instance metric, including the DB instance identifier, timestamp, minimum, maximum, sum, sample count, and unit of measurement for you.

The `aws_rds_db_instance_metric_read_iops` table provides you with metric statistics at 5 minute intervals for the most recent 5 days.

## Examples

### Basic info
Explore the performance metrics of your AWS RDS database instances over time. This can be crucial for identifying trends, optimizing performance, and planning for future capacity needs.

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
  aws_rds_db_instance_metric_read_iops
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
  aws_rds_db_instance_metric_read_iops
order by
  db_instance_identifier,
  timestamp;
```

### Intervals where volumes exceed 1000 average read ops
Explore instances when the average read operations on your AWS RDS DB instances exceed 1000. This could help you identify potential overuse or performance issues and take appropriate action.

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
  aws_rds_db_instance_metric_read_iops
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
  aws_rds_db_instance_metric_read_iops
where
  average > 1000
order by
  db_instance_identifier,
  timestamp;
```

### Intervals where volumes exceed 8000 max read ops
Identify instances where the maximum read operations exceed 8000 in your AWS RDS database instances. This can help in analyzing performance patterns and pinpointing potential areas for optimization.

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
  aws_rds_db_instance_metric_read_iops
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
  aws_rds_db_instance_metric_read_iops
where
  maximum > 8000
order by
  db_instance_identifier,
  timestamp;
```

### Read, Write, and Total IOPS
Gain insights into the performance of your AWS RDS instances by analyzing input/output operations per second (IOPS). This query allows you to monitor the average, maximum, and minimum read and write operations, which can help optimize database performance and capacity planning.

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