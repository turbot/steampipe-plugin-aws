---
title: "Table: aws_rds_db_instance_metric_cpu_utilization_hourly - Query AWS RDS DB Instance Metrics using SQL"
description: "Allows users to query AWS RDS DB Instance CPU Utilization Metrics on an hourly basis."
---

# Table: aws_rds_db_instance_metric_cpu_utilization_hourly - Query AWS RDS DB Instance Metrics using SQL

The `aws_rds_db_instance_metric_cpu_utilization_hourly` table in Steampipe provides information about the CPU utilization metrics of AWS RDS DB instances on an hourly basis. This table allows DevOps engineers to query specific details about CPU usage, including maximum, minimum, and average utilization, as well as the sum of all utilization within the specified time frame. Users can utilize this table to monitor and analyze the CPU consumption of their RDS DB instances, which can aid in optimizing resource usage and identifying potential performance issues. The schema outlines the various attributes of the CPU utilization metric, including the DB instance identifier, timestamp, and various statistics related to CPU utilization.

The `aws_rds_db_instance_metric_cpu_utilization_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_rds_db_instance_metric_cpu_utilization_hourly` table, you can use the `.inspect aws_rds_db_instance_metric_cpu_utilization_hourly` command in Steampipe.

### Key columns:

- `db_instance_identifier`: This is the identifier for the DB instance. It is a key column because it uniquely identifies the DB instance and can be used to join with other tables that contain DB instance information.
- `timestamp`: This column represents the time at which the CPU utilization data was recorded. It is important as it allows for tracking of CPU utilization over time.
- `average`: This column represents the average CPU utilization for the DB instance within the specified time frame. It is useful for understanding the typical CPU usage of the DB instance.

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
  aws_rds_db_instance_metric_cpu_utilization_hourly
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
  aws_rds_db_instance_metric_cpu_utilization_hourly
where average > 80
order by
  db_instance_identifier,
  timestamp;
```

### CPU hourly average < 2%

```sql
select
  db_instance_identifier,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_rds_db_instance_metric_cpu_utilization_hourly
where average < 2
order by
  db_instance_identifier,
  timestamp;
```