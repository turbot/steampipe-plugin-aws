---
title: "Table: aws_redshift_cluster_metric_cpu_utilization_daily - Query AWS Redshift Cluster Metrics using SQL"
description: "Allows users to query AWS Redshift Cluster CPU Utilization Metrics on a daily basis."
---

# Table: aws_redshift_cluster_metric_cpu_utilization_daily - Query AWS Redshift Cluster Metrics using SQL

The `aws_redshift_cluster_metric_cpu_utilization_daily` table in Steampipe provides information about the CPU utilization metrics for AWS Redshift clusters, calculated on a daily basis. This table allows data engineers and administrators to query CPU usage details, including maximum, minimum, average, and sample count values. Users can utilize this table to gather insights on cluster performance, such as identifying clusters with high CPU usage, tracking CPU usage trends over time, and optimizing cluster configurations for better performance. The schema outlines the various attributes of the CPU utilization metrics, including the cluster identifier, region, timestamp, and the aforementioned statistical values.

The `aws_redshift_cluster_metric_cpu_utilization_daily` table provides metric statistics at 24 hour intervals for the last year.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_redshift_cluster_metric_cpu_utilization_daily` table, you can use the `.inspect aws_redshift_cluster_metric_cpu_utilization_daily` command in Steampipe.

### Key columns:

- `cluster_identifier`: The identifier of the cluster for which the CPU utilization metrics are being reported. This column can be used to join this table with other tables that contain cluster-specific information.
- `region`: The AWS region in which the cluster is located. This column can be used to join this table with other tables that contain region-specific information.
- `timestamp`: The timestamp for the reported CPU utilization metrics. This column can be used to join this table with other tables that contain time-series data.

## Examples

### Basic info

```sql
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

```sql
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

### CPU daily average < 2%

```sql
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