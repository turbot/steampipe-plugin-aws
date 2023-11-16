---
title: "Table: aws_rds_db_instance_metric_read_iops_daily - Query AWS RDS DBInstance using SQL"
description: "Allows users to query AWS RDS DBInstance metrics for daily read IOPS (Input/Output Operations Per Second)."
---

# Table: aws_rds_db_instance_metric_read_iops_daily - Query AWS RDS DBInstance using SQL

The `aws_rds_db_instance_metric_read_iops_daily` table in Steampipe provides information about the daily read IOPS metrics of AWS RDS DBInstances. This table allows DevOps engineers to query DBInstance-specific details, including the number of read I/O operations from the DBInstance per day. Users can utilize this table to gather insights on DBInstance performance, such as identifying DBInstances that have a high read I/O operations rate, which could indicate potential performance bottlenecks. The schema outlines the various attributes of the DBInstance's daily read IOPS metrics, including the DBInstance identifier, timestamp of the metric, and the minimum, maximum, and average number of read IOPS.

The `aws_rds_db_instance_metric_read_iops_daily` table provides metric statistics at 24 hour intervals for the last year.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_rds_db_instance_metric_read_iops_daily` table, you can use the `.inspect aws_rds_db_instance_metric_read_iops_daily` command in Steampipe.

### Key columns:

- `db_instance_identifier`: The identifier of the DBInstance. This column is useful for joining with other tables that also contain DBInstance identifiers to get more detailed information about the DBInstance.
- `timestamp`: The timestamp for the specific metric data. This column is important for tracking the DBInstance's read IOPS over time.
- `average`: The average number of read IOPS for the DBInstance during the day. This column is useful for identifying trends in the DBInstance's performance.

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
  aws_rds_db_instance_metric_read_iops_daily
order by
  db_instance_identifier,
  timestamp;
```

### Intervals where volumes exceed 1000 average read ops
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
  aws_rds_db_instance_metric_read_iops_daily
where
  average > 1000
order by
  db_instance_identifier,
  timestamp;
```


### Intervals where volumes exceed 8000 max read ops
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
  aws_rds_db_instance_metric_read_iops_daily
where
  maximum > 8000
order by
  db_instance_identifier,
  timestamp;
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
  aws_rds_db_instance_metric_read_iops_daily as r,
  aws_rds_db_instance_metric_write_iops_daily as w
where 
  r.db_instance_identifier = w.db_instance_identifier
  and r.timestamp = w.timestamp
order by
  r.db_instance_identifier,
  r.timestamp;
```