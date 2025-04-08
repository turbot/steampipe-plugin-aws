---
title: "Steampipe Table: aws_rds_db_instance_metric_connections_daily - Query AWS RDS DB Instance Metrics using SQL"
description: "Allows users to query AWS RDS DB Instance Metrics on a daily basis, retrieving information about the number of database connections."
folder: "RDS"
---

# Table: aws_rds_db_instance_metric_connections_daily - Query AWS RDS DB Instance Metrics using SQL

The AWS RDS DB Instance is a relational database service that provides cost-efficient and resizable capacity while managing time-consuming database administration tasks. It provides a set of metrics, such as the "Connections" metric, to help you monitor the performance of the DB instances. The "Connections" metric, specifically, measures the number of database connections made to an RDS instance on a daily basis, aiding in understanding usage patterns and potential performance issues.

## Table Usage Guide

The `aws_rds_db_instance_metric_connections_daily` table in Steampipe provides you with information about AWS RDS DB instance metrics, specifically focusing on the daily database connections. This table allows you, as a DevOps engineer, database administrator, or other technical professional, to query these metrics, enabling you to monitor and manage the number of connections to your database instances. This is critical for your optimization of resource utilization, management of performance, and troubleshooting of issues. The schema outlines various attributes of the daily DB instance connections metric, including the DB instance identifier, timestamp, sum of connections, minimum, maximum, and sample count.

The `aws_rds_db_instance_metric_connections_daily` table provides you with metric statistics at 24 hour intervals for the past year.

## Examples

### Basic info
Analyze the daily connection metrics of your AWS RDS database instances to understand their usage patterns and performance. This can help in identifying bottlenecks, planning capacity, and optimizing resource utilization.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which your AWS RDS database instances have an average of over 100 daily connections. This can help in understanding the load on your databases and potentially optimize them for better performance.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which your database instances have not had any connections in the past week. This can help you identify unused or idle instances, potentially saving on unnecessary costs.

```sql+postgres
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

```sql+sqlite
select
  db_instance_identifier,
  sum(maximum) as total_connections
from
  aws_rds_db_instance_metric_connections
where 
  timestamp > (date('now', '-7 day'))
group by
  db_instance_identifier
having
  sum(maximum) = 0 
;
```