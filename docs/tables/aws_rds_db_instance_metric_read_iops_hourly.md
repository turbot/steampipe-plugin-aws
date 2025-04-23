---
title: "Steampipe Table: aws_rds_db_instance_metric_read_iops_hourly - Query AWS RDS DB Instances using SQL"
description: "Allows users to query AWS RDS DB Instances and retrieve hourly metrics related to read IOPS (Input/Output Operations Per Second)."
folder: "RDS"
---

# Table: aws_rds_db_instance_metric_read_iops_hourly - Query AWS RDS DB Instances using SQL

The AWS RDS DB Instance is a part of Amazon's Relational Database Service that provides resizable capacity for an industry-standard relational database and manages common database administration tasks. It offers high availability and security, along with compatibility with several database engines, including MySQL, Oracle, and PostgreSQL. The 'aws_rds_db_instance_metric_read_iops_hourly' represents the Input/Output Operations Per Second (IOPS) for read operations on the DB instance, measured on an hourly basis.

## Table Usage Guide

The `aws_rds_db_instance_metric_read_iops_hourly` table in Steampipe provides you with information about the read IOPS metrics for AWS Relational Database Service (RDS) DB instances on an hourly basis. This table enables you, as a DevOps engineer, database administrator, or other technical professional, to query read IOPS metrics. These can be useful for monitoring database performance, planning for capacity, and troubleshooting performance issues. The schema outlines various attributes of the read IOPS metrics for you, including the timestamp, average, maximum, and minimum read IOPS, as well as the standard deviation.

The `aws_rds_db_instance_metric_read_iops_hourly` table provides you with metric statistics at 1 hour intervals for the most recent 60 days.

## Examples

### Basic info
Analyze the performance of your AWS RDS database instances over time to optimize resource allocation and improve efficiency. This query helps you understand the input/output operations per second (IOPS) on an hourly basis, allowing for effective capacity planning and performance tuning.

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
  aws_rds_db_instance_metric_read_iops_hourly
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
  aws_rds_db_instance_metric_read_iops_hourly
order by
  db_instance_identifier,
  timestamp;
```

### Intervals where volumes exceed 1000 average read ops
Identify instances where the average read operations per second on your Amazon RDS database instances exceed 1000 within an hour. This can help you manage and optimize your database performance by pinpointing periods of high demand.

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
  aws_rds_db_instance_metric_read_iops_hourly
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
  aws_rds_db_instance_metric_read_iops_hourly
where
  average > 1000
order by
  db_instance_identifier,
  timestamp;
```


### Intervals where volumes exceed 8000 max read ops
Discover the instances when the read operations on your AWS RDS database instances exceed a certain threshold, in this case, 8000 maximum read operations per hour. This can be particularly useful for identifying periods of high load or potential performance issues.

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
  aws_rds_db_instance_metric_read_iops_hourly
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
  aws_rds_db_instance_metric_read_iops_hourly
where
  maximum > 8000
order by
  db_instance_identifier,
  timestamp;
```



### Intervals where volume average iops exceeds provisioned iops
Determine the instances where the average input/output operations per second (IOPS) exceeds the provisioned IOPS for your Amazon RDS database instances. This can help you identify periods of high load, enabling you to better plan for capacity and avoid potential performance issues.

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
Analyze the input/output operations per second (IOPS) for your AWS RDS database instances to understand their read and write performance over time. This can help in optimizing resources and planning for capacity needs.

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