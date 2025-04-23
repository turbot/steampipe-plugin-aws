---
title: "Steampipe Table: aws_rds_db_instance_metric_connections - Query AWS RDS DBInstance Metrics using SQL"
description: "Allows users to query AWS RDS DBInstance Metrics for a comprehensive view of the number of database connections."
folder: "RDS"
---

# Table: aws_rds_db_instance_metric_connections - Query AWS RDS DBInstance Metrics using SQL

The AWS RDS DBInstance Metrics is a feature of Amazon Relational Database Service (RDS) that allows you to monitor the performance of your databases. It provides a variety of metrics that can help you understand your database's workload, throughput, and performance. These metrics can be queried using SQL, enabling you to analyze and optimize your database's performance.

## Table Usage Guide

The `aws_rds_db_instance_metric_connections` table in Steampipe provides you with information about the number of database connections to each Amazon RDS DB instance. This table allows you, as a DevOps engineer, database administrator, or other technical professional, to query connection-specific details, including the time at which the number of connections was recorded, the maximum number of connections during the period, and the number of data points used for the statistical calculation. You can utilize this table to monitor and manage your database connections, analyze connection trends, and troubleshoot potential connection issues. The schema outlines the various attributes of the RDS DB instance connections, including the DB instance identifier, timestamp, sample count, average, minimum, and maximum number of connections.

The `aws_rds_db_instance_metric_connections` table provides you with metric statistics at 5 minute intervals for the most recent 5 days.

## Examples

### Basic info
Determine the areas in which AWS RDS database instances have varying connection metrics over time. This can help in understanding the database's performance and planning for scalability.

```sql+postgres
select
  db_instance_identifier,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_rds_db_instance_metric_connections
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
  aws_rds_db_instance_metric_connections
order by
  db_instance_identifier,
  timestamp;
```



### Intervals averaging over 100 connections
Determine the areas in which your AWS RDS database instances have an average of over 100 connections, allowing you to identify potential performance issues or heavy usage periods.

```sql+postgres
select
  db_instance_identifier,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_rds_db_instance_metric_connections
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
  aws_rds_db_instance_metric_connections
where 
  average > 100
order by
  db_instance_identifier,
  timestamp;
```