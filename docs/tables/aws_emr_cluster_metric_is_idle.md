---
title: "Steampipe Table: aws_emr_cluster_metric_is_idle - Query AWS EMR Cluster Metrics using SQL"
description: "Allows users to query AWS EMR Cluster Metrics to determine if a cluster is idle."
folder: "EMR"
---

# Table: aws_emr_cluster_metric_is_idle - Query AWS EMR Cluster Metrics using SQL

The AWS Elastic MapReduce (EMR) Cluster is a managed service that simplifies big data processing, with capabilities to analyze and process vast amounts of data quickly. It utilizes popular distributed frameworks such as Apache Hadoop and Spark. The 'is_idle' metric indicates whether the cluster is idle or not, helping to optimize resource utilization and cost.

## Table Usage Guide

The `aws_emr_cluster_metric_is_idle` table in Steampipe provides you with information about the IsIdle metric of AWS Elastic MapReduce (EMR) Clusters. This table enables you, as a DevOps engineer, to query details related to the idle state of EMR clusters, including the cluster's ID, the namespace of the metric, and the timestamp of the metric. You can utilize this table to gather insights on cluster utilization, such as identifying underused resources and optimizing resource allocation. The schema outlines the various attributes of the EMR Cluster IsIdle metric for you, including the cluster ID, metric namespace, and metric timestamp.

The `aws_emr_cluster_metric_is_idle` table provides you with metric statistics at 5 minute intervals for the most recent 5 days.

## Examples

### Basic info
Assess the performance of your AWS EMR clusters by analyzing their idle metrics over time. This helps in optimizing resource usage by identifying periods of low activity or inactivity.

```sql+postgres
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_emr_cluster_metric_is_idle
order by
  id,
  timestamp;
```

```sql+sqlite
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_emr_cluster_metric_is_idle
order by
  id,
  timestamp;
```