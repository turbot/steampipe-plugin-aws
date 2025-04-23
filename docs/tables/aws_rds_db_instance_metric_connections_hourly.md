---
title: "Steampipe Table: aws_rds_db_instance_metric_connections_hourly - Query AWS RDS DB Instance Metrics using SQL"
description: "Allows users to query AWS RDS DB Instance Metrics on an hourly basis, specifically the connection metrics. It provides data about the number of database connections to each DB instance in your Amazon RDS environment."
folder: "RDS"
---

# Table: aws_rds_db_instance_metric_connections_hourly - Query AWS RDS DB Instance Metrics using SQL

The AWS RDS DB Instance Metrics is a service that allows you to monitor database instances. It provides key metrics related to connections, which are computed on an hourly basis, helping you understand the usage and performance of your databases. These metrics can be queried using SQL, enabling easy integration with existing monitoring systems or custom analysis.

## Table Usage Guide

The `aws_rds_db_instance_metric_connections_hourly` table in Steampipe provides you with information about the connection metrics for each DB instance in your Amazon RDS environment, aggregated on an hourly basis. This table allows you, as a DevOps engineer, database administrator, or other technical professional, to query connection-specific details, including the number of connections, the time of the connections, and associated metadata. You can utilize this table to gather insights on connection patterns, such as peak connection times, connection durations, and more. The schema outlines the various attributes of the DB instance connection metrics for you, including the timestamp, maximum and minimum number of connections, and the sample count.

The `aws_rds_db_instance_metric_connections_hourly` table provides you with metric statistics at 1 hour intervals for the most recent 60 days.

## Examples

### Basic info
Explore the performance of your AWS RDS instances by examining the minimum, maximum, and average hourly connection metrics. This allows you to identify patterns and potential issues, ensuring optimal performance and resource allocation.

```sql+postgres
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

```sql+sqlite
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
This example helps you identify instances in your AWS RDS database where the average number of connections exceeds 100 in an hour. It's useful for monitoring heavy traffic periods and potential performance issues in your database.

```sql+postgres
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

```sql+sqlite
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