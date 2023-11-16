---
title: "Table: aws_rds_db_instance_metric_cpu_utilization_daily - Query AWS RDS DB Instances using SQL"
description: "Allows users to query AWS RDS DB Instances to retrieve daily CPU utilization metrics."
---

# Table: aws_rds_db_instance_metric_cpu_utilization_daily - Query AWS RDS DB Instances using SQL

The `aws_rds_db_instance_metric_cpu_utilization_daily` table in Steampipe provides information about the daily CPU utilization metrics of AWS RDS DB Instances. This table allows DevOps engineers, database administrators, and other technical professionals to query CPU-specific details, including maximum and average CPU utilization, and timestamps. Users can utilize this table to monitor and analyze the CPU usage patterns of RDS DB Instances over time. The schema outlines the various attributes of the CPU utilization metrics, including the DB instance identifier, timestamp, maximum utilization, and average utilization.

The `aws_rds_db_instance_metric_cpu_utilization_daily` table provides metric statistics at 24 hour intervals for the last year.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_rds_db_instance_metric_cpu_utilization_daily` table, you can use the `.inspect aws_rds_db_instance_metric_cpu_utilization_daily` command in Steampipe.

**Key columns**:

- `db_instance_identifier`: The identifier of the DB instance. This column is useful as it can be used to join this table with other tables to get more detailed information about the DB instance.
- `timestamp`: The timestamp of the data point for the CPU utilization metric. This column is important as it allows users to track CPU utilization over time.
- `maximum`: The maximum CPU utilization for the DB instance for the day. This column is useful for identifying peak CPU usage periods.

## Examples

### Basic info

```sql
select
  db_instance_identifier,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_rds_db_instance_metric_cpu_utilization_daily
order by
  db_instance_identifier,
  timestamp;
```



### CPU Over 80% average

```sql
select
  db_instance_identifier,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_rds_db_instance_metric_cpu_utilization_daily
where average > 80
order by
  db_instance_identifier,
  timestamp;
```

### CPU daily average < 2%

```sql
select
  db_instance_identifier,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_rds_db_instance_metric_cpu_utilization_daily
where average < 2
order by
  db_instance_identifier,
  timestamp;
```