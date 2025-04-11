---
title: "Steampipe Table: aws_rds_db_instance_metric_cpu_utilization_daily - Query AWS RDS DB Instances using SQL"
description: "Allows users to query AWS RDS DB Instances to retrieve daily CPU utilization metrics."
folder: "RDS"
---

# Table: aws_rds_db_instance_metric_cpu_utilization_daily - Query AWS RDS DB Instances using SQL

The AWS RDS DB Instance is a part of Amazon Relational Database Service (RDS), a web service that makes it easier to set up, operate, and scale a relational database in the cloud. It provides cost-efficient, resizable capacity for an industry-standard relational database and manages common database administration tasks. The 'cpu_utilization_daily' metric provides the percentage of CPU utilization for an Amazon RDS instance, averaged over a 24-hour period.

## Table Usage Guide

The `aws_rds_db_instance_metric_cpu_utilization_daily` table in Steampipe provides you with information about the daily CPU utilization metrics of AWS RDS DB Instances. This table allows you, as a DevOps engineer, database administrator, or other technical professional, to query CPU-specific details, including maximum and average CPU utilization, and timestamps. You can utilize this table to monitor and analyze the CPU usage patterns of RDS DB Instances over time. The schema outlines the various attributes of the CPU utilization metrics for you, including the DB instance identifier, timestamp, maximum utilization, and average utilization.

The `aws_rds_db_instance_metric_cpu_utilization_daily` table provides you with metric statistics at 24 hour intervals for the last year.

## Examples

### Basic info
Analyze the daily CPU utilization of your AWS RDS database instances to understand their performance trends and identify any instances that may be under or over-utilized. This will assist in optimizing resource allocation and potentially reducing costs.

```sql+postgres
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

```sql+sqlite
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
Explore which AWS RDS database instances have an average CPU utilization exceeding 80%, allowing you to proactively manage and optimize your resources for better performance and cost efficiency.

```sql+postgres
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

```sql+sqlite
select
  db_instance_identifier,
  timestamp,
  round(minimum,2) as min_cpu,
  round(maximum,2) as max_cpu,
  round(average,2) as avg_cpu,
  sample_count
from
  aws_rds_db_instance_metric_cpu_utilization_daily
where average > 80
order by
  db_instance_identifier,
  timestamp;
```

### CPU daily average < 2%
Explore which database instances have a daily average CPU utilization less than 2%. This can help in identifying underutilized resources and potentially save costs.

```sql+postgres
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

```sql+sqlite
select
  db_instance_identifier,
  timestamp,
  round(minimum,2) as min_cpu,
  round(maximum,2) as max_cpu,
  round(average,2) as avg_cpu,
  sample_count
from
  aws_rds_db_instance_metric_cpu_utilization_daily
where average < 2
order by
  db_instance_identifier,
  timestamp;
```