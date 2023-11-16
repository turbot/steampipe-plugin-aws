---
title: "Table: aws_rds_db_instance_metric_write_iops - Query AWS RDS DBInstance Write IOPS using SQL"
description: "Allows users to query AWS RDS DBInstance Write IOPS to retrieve metrics on the write input/output operations per second."
---

# Table: aws_rds_db_instance_metric_write_iops - Query AWS RDS DBInstance Write IOPS using SQL

The `aws_rds_db_instance_metric_write_iops` table in Steampipe provides information about the Write IOPS (Input/Output Operations Per Second) metrics of AWS RDS DBInstances. This table allows DevOps engineers to query details related to the write operations performance of their RDS DBInstances, including the average, maximum, and minimum values for a specified period. Users can utilize this table to monitor and analyze the write performance of their DBInstances, helping them optimize the performance and reliability of their database operations. The schema outlines the various attributes of the Write IOPS metrics, including the DBInstance identifier, timestamp, and the statistics for the period.

The `aws_rds_db_instance_metric_write_iops` table provides metric statistics at 5 minute intervals for the most recent 5 days.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_rds_db_instance_metric_write_iops` table, you can use the `.inspect aws_rds_db_instance_metric_write_iops` command in Steampipe.

### Key columns:

* `db_instance_identifier`: The identifier of the DBInstance. This column can be used to join with the `aws_rds_db_instance` table to get more information about the DBInstance.

* `timestamp`: The timestamp for the data point in UTC. This column is useful for tracking the write performance over time.

* `average`: The average of metric values that correspond to the specified time period. This column provides insight into the average write performance during the period.

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
  aws_rds_db_instance_metric_write_iops
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
  aws_rds_db_instance_metric_write_iops
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
  aws_rds_db_instance_metric_write_iops
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
  aws_rds_db_instance_metric_read_iops as r,
  aws_rds_db_instance_metric_write_iops as w
where 
  r.db_instance_identifier = w.db_instance_identifier
  and r.timestamp = w.timestamp
order by
  r.db_instance_identifier,
  r.timestamp;
```