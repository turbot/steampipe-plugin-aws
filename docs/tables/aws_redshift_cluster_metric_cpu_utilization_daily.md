---
title: "Steampipe Table: aws_redshift_cluster_metric_cpu_utilization_daily - Query AWS Redshift Cluster Metrics using SQL"
description: "Allows users to query AWS Redshift Cluster CPU Utilization Metrics on a daily basis."
folder: "Redshift"
---

# Table: aws_redshift_cluster_metric_cpu_utilization_daily - Query AWS Redshift Cluster Metrics using SQL

The AWS Redshift Cluster is a fully managed, petabyte-scale data warehouse service in the cloud. It allows you to analyze all your data using your existing business intelligence tools. The CPU Utilization metric provides the percentage of CPU utilization for the Amazon Redshift cluster.

## Table Usage Guide

The `aws_redshift_cluster_metric_cpu_utilization_daily` table in Steampipe gives you information about the CPU utilization metrics for AWS Redshift clusters, calculated on a daily basis. This table allows you, as a data engineer or administrator, to query CPU usage details, including maximum, minimum, average, and sample count values. You can utilize this table to gather insights on cluster performance, such as identifying clusters with high CPU usage, tracking CPU usage trends over time, and optimizing cluster configurations for better performance. The schema outlines the various attributes of the CPU utilization metrics, including the cluster identifier, region, timestamp, and the aforementioned statistical values.

The `aws_redshift_cluster_metric_cpu_utilization_daily` table provides you with metric statistics at 24-hour intervals for the last year.

## Examples

### Basic info
Analyze the daily CPU utilization patterns of your AWS Redshift clusters to understand their performance and resource usage trends. This information can help optimize resource allocation and improve overall system efficiency.

```sql+postgres
select
  cluster_identifier,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_redshift_cluster_metric_cpu_utilization_daily
order by
  cluster_identifier,
  timestamp;
```

```sql+sqlite
select
  cluster_identifier,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_redshift_cluster_metric_cpu_utilization_daily
order by
  cluster_identifier,
  timestamp;
```

### CPU Over 80% average
Determine the areas in which your AWS Redshift clusters are experiencing high CPU utilization, specifically where the average daily usage exceeds 80%. This can help in identifying potential performance issues and planning for capacity upgrades.

```sql+postgres
select
  cluster_identifier,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_redshift_cluster_metric_cpu_utilization_daily
where average > 80
order by
  cluster_identifier,
  timestamp;
```

```sql+sqlite
select
  cluster_identifier,
  timestamp,
  round(minimum,2) as min_cpu,
  round(maximum,2) as max_cpu,
  round(average,2) as avg_cpu,
  sample_count
from
  aws_redshift_cluster_metric_cpu_utilization_daily
where average > 80
order by
  cluster_identifier,
  timestamp;
```

### CPU daily average < 2%
This example helps to identify instances where the average daily CPU utilization is less than 2% in your AWS Redshift clusters. This can be useful to pinpoint underutilized resources, potentially leading to cost savings by downsizing or eliminating these clusters.

```sql+postgres
select
  cluster_identifier,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_redshift_cluster_metric_cpu_utilization_daily
where average < 2
order by
  cluster_identifier,
  timestamp;
```

```sql+sqlite
select
  cluster_identifier,
  timestamp,
  round(minimum,2) as min_cpu,
  round(maximum,2) as max_cpu,
  round(average,2) as avg_cpu,
  sample_count
from
  aws_redshift_cluster_metric_cpu_utilization_daily
where average < 2
order by
  cluster_identifier,
  timestamp;
```