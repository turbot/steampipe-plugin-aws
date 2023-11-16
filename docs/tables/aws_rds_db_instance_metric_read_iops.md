---
title: "Table: aws_rds_db_instance_metric_read_iops - Query AWS RDS DBInstanceMetricReadIops using SQL"
description: "Allows users to query AWS RDS DBInstanceMetricReadIops to retrieve and monitor the read IOPS (Input/Output Operations Per Second) metrics for Amazon RDS DB instances."
---

# Table: aws_rds_db_instance_metric_read_iops - Query AWS RDS DBInstanceMetricReadIops using SQL

The `aws_rds_db_instance_metric_read_iops` table in Steampipe provides information about the read IOPS metrics of AWS RDS DB instances. This table allows users such as DevOps engineers and database administrators to query and monitor the read IOPS metrics, which can be useful for performance tuning and capacity planning. The read IOPS refers to the number of read input/output operations per second. The schema outlines the various attributes of the DB instance metric, including the DB instance identifier, timestamp, minimum, maximum, sum, sample count, and unit of measurement.

The `aws_rds_db_instance_metric_read_iops` table provides metric statistics at 5 minute intervals for the most recent 5 days.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_rds_db_instance_metric_read_iops` table, you can use the `.inspect aws_rds_db_instance_metric_read_iops` command in Steampipe.

**Key columns**:

- `db_instance_identifier`: The identifier of the DB instance. This can be used to join this table with other tables that contain information about the DB instance.
- `timestamp`: The timestamp for the data point. This can be used to track the read IOPS over time.
- `average`: The average of metric values that correspond to the data point. This can be used to monitor the average read IOPS.

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
  aws_rds_db_instance_metric_read_iops
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
  aws_rds_db_instance_metric_read_iops
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
  aws_rds_db_instance_metric_read_iops
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