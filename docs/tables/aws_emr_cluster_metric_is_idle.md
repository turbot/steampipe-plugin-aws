---
title: "Table: aws_emr_cluster_metric_is_idle - Query AWS EMR Cluster Metrics using SQL"
description: "Allows users to query AWS EMR Cluster Metrics to determine if a cluster is idle."
---

# Table: aws_emr_cluster_metric_is_idle - Query AWS EMR Cluster Metrics using SQL

The `aws_emr_cluster_metric_is_idle` table in Steampipe provides information about the IsIdle metric of AWS Elastic MapReduce (EMR) Clusters. This table allows DevOps engineers to query details related to the idle state of EMR clusters, including the cluster's ID, the namespace of the metric, and the timestamp of the metric. Users can utilize this table to gather insights on cluster utilization, such as identifying underused resources and optimizing resource allocation. The schema outlines the various attributes of the EMR Cluster IsIdle metric, including the cluster ID, metric namespace, and metric timestamp.

The `aws_emr_cluster_metric_is_idle` table provides metric statistics at 5 minute intervals for the most recent 5 days.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_emr_cluster_metric_is_idle` table, you can use the `.inspect aws_emr_cluster_metric_is_idle` command in Steampipe.

Key columns:

- `cluster_id`: The ID of the EMR cluster. It can be used to join this table with other EMR cluster tables to get more comprehensive information about the cluster.
- `namespace`: The namespace of the AWS CloudWatch metric. This column is useful as it can help in filtering and categorizing the metrics.
- `timestamp`: The timestamp of the metric. This column is useful in tracking the history and changes in the idle state of the cluster over time.

## Examples

### Basic info

```sql
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
