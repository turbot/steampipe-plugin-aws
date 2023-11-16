---
title: "Table: aws_rds_db_instance_metric_connections_daily - Query AWS RDS DB Instance Metrics using SQL"
description: "Allows users to query AWS RDS DB Instance Metrics on a daily basis, retrieving information about the number of database connections."
---

# Table: aws_rds_db_instance_metric_connections_daily - Query AWS RDS DB Instance Metrics using SQL

The `aws_rds_db_instance_metric_connections_daily` table in Steampipe provides information about AWS RDS DB instance metrics, specifically focusing on the daily database connections. This table allows DevOps engineers, database administrators, and other technical professionals to query these metrics, enabling them to monitor and manage the number of connections to their database instances. This is critical for optimizing resource utilization, managing performance, and troubleshooting issues. The schema outlines various attributes of the daily DB instance connections metric, including the DB instance identifier, timestamp, sum of connections, minimum, maximum, and sample count.

The `aws_rds_db_instance_metric_connections_daily` table provides metric statistics at 24 hour intervals for the past year.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_rds_db_instance_metric_connections_daily` table, you can use the `.inspect aws_rds_db_instance_metric_connections_daily` command in Steampipe.

**Key columns**:

- `db_instance_identifier`: This is the identifier for the DB instance. It is crucial for joining this table with other tables that contain DB instance-specific information.
- `timestamp`: This column holds the date and time of the metric data point. It is important for time-series analysis and tracking changes over time.
- `sum`: This column represents the sum of the metric values for the specified period. It is useful for understanding the total number of connections made to the DB instance within a day.

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
  aws_rds_db_instance_metric_connections_daily
order by
  db_instance_identifier,
  timestamp;
```


### Intervals averaging over 100 connections

```sql
select
  db_instance_identifier,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_rds_db_instance_metric_connections_daily
where 
  average > 100
order by
  db_instance_identifier,
  timestamp;
```


### Instances with no connections in the past week

```sql
select
  db_instance_identifier,
  sum(maximum) as total_connections
from
  aws_rds_db_instance_metric_connections
where 
  timestamp > (current_date - interval '7' day)
group by
  db_instance_identifier
having
  sum(maximum) = 0 
;
```


