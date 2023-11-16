---
title: "Table: aws_rds_db_instance_metric_read_iops_hourly - Query AWS RDS DB Instances using SQL"
description: "Allows users to query AWS RDS DB Instances and retrieve hourly metrics related to read IOPS (Input/Output Operations Per Second)."
---

# Table: aws_rds_db_instance_metric_read_iops_hourly - Query AWS RDS DB Instances using SQL

The `aws_rds_db_instance_metric_read_iops_hourly` table in Steampipe provides information about the read IOPS metrics for AWS Relational Database Service (RDS) DB instances on an hourly basis. This table allows DevOps engineers, database administrators, and other technical professionals to query read IOPS metrics, which can be useful for monitoring database performance, planning for capacity, and troubleshooting performance issues. The schema outlines various attributes of the read IOPS metrics, including the timestamp, average, maximum, and minimum read IOPS, as well as the standard deviation.

The `aws_rds_db_instance_metric_read_iops_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_rds_db_instance_metric_read_iops_hourly` table, you can use the `.inspect aws_rds_db_instance_metric_read_iops_hourly` command in Steampipe.

**Key columns**:

- `db_instance_identifier`: The identifier for the DB instance. This column is important as it uniquely identifies the DB instance for which the read IOPS metrics are being retrieved. It can be used to join this table with other tables that contain information about the DB instance.

- `timestamp`: The date and time at which the read IOPS metrics were recorded. This column is useful for tracking the performance of the DB instance over time.

- `average`: The average read IOPS for the hour. This column is significant as it gives an overview of the DB instance's read performance during the hour. It can be useful for identifying trends and patterns in the DB instance's performance.

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
  aws_rds_db_instance_metric_read_iops_hourly
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
  aws_rds_db_instance_metric_read_iops_hourly
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
  aws_rds_db_instance_metric_read_iops_hourly
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