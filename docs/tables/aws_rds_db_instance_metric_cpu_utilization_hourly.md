---
title: "Steampipe Table: aws_rds_db_instance_metric_cpu_utilization_hourly - Query AWS RDS DB Instance Metrics using SQL"
description: "Allows users to query AWS RDS DB Instance CPU Utilization Metrics on an hourly basis."
folder: "RDS"
---

# Table: aws_rds_db_instance_metric_cpu_utilization_hourly - Query AWS RDS DB Instance Metrics using SQL

The AWS RDS DB Instance Metrics is a feature of Amazon Relational Database Service (RDS) that allows you to monitor the performance of your databases. It provides CPU utilization metrics, which indicate the percentage of CPU utilization for an Amazon RDS instance. This can be used to assess the load on your database and optimize performance as necessary.

## Table Usage Guide

The `aws_rds_db_instance_metric_cpu_utilization_hourly` table in Steampipe provides you with information about the CPU utilization metrics of AWS RDS DB instances on an hourly basis. This table enables you, as a DevOps engineer, to query specific details about CPU usage, including maximum, minimum, and average utilization, as well as the sum of all utilization within the specified time frame. You can utilize this table to monitor and analyze the CPU consumption of your RDS DB instances, which can assist you in optimizing resource usage and identifying potential performance issues. The schema outlines the various attributes of the CPU utilization metric for you, including the DB instance identifier, timestamp, and various statistics related to CPU utilization.

The `aws_rds_db_instance_metric_cpu_utilization_hourly` table provides you with metric statistics at 1 hour intervals for the most recent 60 days.

## Examples


### Basic info
Analyze the CPU utilization of AWS RDS database instances over time to understand performance trends and identify potential bottlenecks or periods of high demand.

```sql+postgres
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

```sql+sqlite
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
Discover the instances where the average CPU utilization of your AWS RDS database instances exceeds 80%. This could be used to identify potential performance issues and manage resources more effectively.

```sql+postgres
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

```sql+sqlite
select
  db_instance_identifier,
  timestamp,
  round(minimum,2) as min_cpu,
  round(maximum,2) as max_cpu,
  round(average,2) as avg_cpu,
  sample_count
from
  aws_rds_db_instance_metric_cpu_utilization_hourly
where average > 80
order by
  db_instance_identifier,
  timestamp;
```

### CPU hourly average < 2%
Explore which AWS RDS database instances have an average CPU utilization of less than 2% on an hourly basis. This can be useful to identify potentially under-utilized resources and optimize infrastructure costs.

```sql+postgres
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

```sql+sqlite
select
  db_instance_identifier,
  timestamp,
  round(minimum,2) as min_cpu,
  round(maximum,2) as max_cpu,
  round(average,2) as avg_cpu,
  sample_count
from
  aws_rds_db_instance_metric_cpu_utilization_hourly
where average < 2
order by
  db_instance_identifier,
  timestamp;
```