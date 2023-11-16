---
title: "Table: aws_rds_db_instance_metric_connections_hourly - Query AWS RDS DB Instance Metrics using SQL"
description: "Allows users to query AWS RDS DB Instance Metrics on an hourly basis, specifically the connection metrics. It provides data about the number of database connections to each DB instance in your Amazon RDS environment."
---

# Table: aws_rds_db_instance_metric_connections_hourly - Query AWS RDS DB Instance Metrics using SQL

The `aws_rds_db_instance_metric_connections_hourly` table in Steampipe provides information about the connection metrics for each DB instance in your Amazon RDS environment, aggregated on an hourly basis. This table allows DevOps engineers, database administrators, and other technical professionals to query connection-specific details, including the number of connections, the time of the connections, and associated metadata. Users can utilize this table to gather insights on connection patterns, such as peak connection times, connection durations, and more. The schema outlines the various attributes of the DB instance connection metrics, including the timestamp, maximum and minimum number of connections, and the sample count.

The `aws_rds_db_instance_metric_connections_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_rds_db_instance_metric_connections_hourly` table, you can use the `.inspect aws_rds_db_instance_metric_connections_hourly` command in Steampipe.

### Key columns:

- `db_instance_identifier`: The identifier for the DB instance. This column can be used to join with other tables that provide more detailed information about each DB instance.
- `timestamp`: The time when the metric data was recorded. This column is useful for tracking connection patterns over time.
- `maximum`: The maximum number of connections recorded during the specified period. This column can be used to identify peak connection times.

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
  aws_rds_db_instance_metric_connections_hourly
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
  aws_rds_db_instance_metric_connections_hourly
where 
  average > 100
order by
  db_instance_identifier,
  timestamp;
```