---
title: "Table: aws_rds_db_instance_metric_write_iops_daily - Query AWS RDS DBInstance using SQL"
description: "Allows users to query AWS RDS DBInstance metrics for daily write IOPS."
---

# Table: aws_rds_db_instance_metric_write_iops_daily - Query AWS RDS DBInstance using SQL

The `aws_rds_db_instance_metric_write_iops_daily` table in Steampipe provides information about the daily write IOPS (Input/Output Operations Per Second) metrics for each AWS RDS DBInstance. This table allows DevOps engineers, DBAs, and other technical professionals to query and analyze the daily write IOPS metrics, which can be critical for performance tuning, capacity planning, and cost management. The schema outlines the various attributes of the daily write IOPS metrics, including the DBInstance identifier, timestamp, minimum, maximum, sum, and average values, among others.

The `aws_rds_db_instance_metric_write_iops_daily` table provides metric statistics at 24 hour intervals for the last year.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_rds_db_instance_metric_write_iops_daily` table, you can use the `.inspect aws_rds_db_instance_metric_write_iops_daily` command in Steampipe.

Key columns:

- `db_instance_identifier`: The identifier of the DBInstance. This column is critical as it uniquely identifies the DBInstance and can be used to join this table with other RDS DBInstance tables.
- `timestamp`: The timestamp of the metric data point. This column is essential for analyzing the daily write IOPS metrics over time.
- `average`: The average of the metric value for the data points during the day. This column is useful for understanding the typical daily write IOPS for the DBInstance.

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
  aws_rds_db_instance_metric_write_iops_daily
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
  aws_rds_db_instance_metric_write_iops_daily
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
  aws_rds_db_instance_metric_write_iops_daily
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