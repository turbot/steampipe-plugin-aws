---
title: "Table: aws_rds_db_instance_metric_write_iops_hourly - Query AWS RDS DBInstance Metrics using SQL"
description: "Allows users to query AWS RDS DBInstance write IOPS metrics on an hourly basis."
---

# Table: aws_rds_db_instance_metric_write_iops_hourly - Query AWS RDS DBInstance Metrics using SQL

The `aws_rds_db_instance_metric_write_iops_hourly` table in Steampipe provides information about the Input/Output operations per second (IOPS) for write operations on an AWS RDS DBInstance, aggregated on an hourly basis. This table allows DevOps engineers, database administrators, and other technical professionals to query DBInstance-specific details, including the number of write IOPS, the timestamp of the data point, and the statistical value. Users can utilize this table to gather insights on the write performance of their DBInstances, such as identifying periods of high write activity, monitoring the impact of performance tuning measures, and more. The schema outlines the various attributes of the DBInstance metric, including the DBInstance identifier, the period, the unit, and the timestamp.

The `aws_rds_db_instance_metric_write_iops_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the aws_rds_db_instance_metric_write_iops_hourly table, you can use the `.inspect aws_rds_db_instance_metric_write_iops_hourly` command in Steampipe.

Key columns:

- `db_instance_identifier`: The identifier of the DBInstance. This can be used to join this table with other tables that also include DBInstance identifiers.
- `timestamp`: The timestamp for the data point. This is useful for tracking changes over time and identifying periods of high write activity.
- `statistic`: The statistical value for the data point. This can be used to analyze the write performance of the DBInstance.

## Examples

### Basic info

```sql
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
```sql
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
```sql
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
```sql
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


### Read, Write, and Total IOPS

```sql
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