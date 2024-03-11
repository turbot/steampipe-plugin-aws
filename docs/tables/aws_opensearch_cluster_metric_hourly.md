---
title: "Steampipe Table: aws_opensearch_cluster_metric_hourly - Query AWS OpenSearch Cluster Metrics using SQL"
description: "Enables users to query AWS OpenSearch Cluster Metrics to obtain detailed insights into cluster performance metrics."
---

# Table: aws_opensearch_cluster_metric_hourly - Query AWS OpenSearch Cluster Metrics using SQL

Amazon OpenSearch Service (successor to Amazon Elasticsearch Service) provides comprehensive metrics that allow for the monitoring of cluster performance. These metrics are crucial for understanding the health and operational status of your OpenSearch clusters. By analyzing these metrics, you can make informed decisions regarding optimization, scaling, and troubleshooting of your clusters over a 24-hour period.

## Table Usage Guide

The `aws_opensearch_cluster_metric_hourly` table in Steampipe allows you to query metric statistics related to the performance and usage of Amazon OpenSearch clusters. This table is a valuable resource for DevOps engineers, cloud architects, and database administrators who need to monitor OpenSearch cluster metrics closely. It provides key metrics such as average, minimum, maximum values, and the sum of metric values over specified periods, along with the number of samples contributing to each metric. By utilizing this table, you can track various metrics over time to identify trends, detect anomalies, and ensure your clusters are running efficiently.

The `aws_opensearch_cluster_metric_hourly` table provides you with metric statistics at 1 hour intervals for the most recent 60 days.

**Important Notes**
- In order to get metric details, the `metric_name` column must be specified. For more information on available metrics for clusters, please see [Amazon OpenSearch Service provides the following metrics for clusters](https://docs.aws.amazon.com/opensearch-service/latest/developerguide/managedomains-cloudwatchmetrics.html#managedomains-cloudwatchmetrics-cluster-metrics).

## Examples

### Basic Cluster CPUUtilization Metrics Overview
Gain a basic understanding of your AWS OpenSearch cluster's metrics, including average, minimum, maximum values, and sample count. This overview helps in monitoring the general health and performance of your clusters.

```sql+postgres
select
  domain_name,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_opensearch_cluster_metric_hourly
where
  metric_name = 'CPUUtilization'
order by
  domain_name,
  timestamp;
```

```sql+sqlite
select
  domain_name,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_opensearch_cluster_metric_hourly
where
  metric_name = 'CPUUtilization'
order by
  domain_name,
  timestamp;
```

### Identifying High Usage Periods
Locate instances where specific metrics, such as CPU utilization or JVM memory pressure, exceed expected thresholds. This analysis helps in pinpointing periods of high usage or stress on your OpenSearch clusters.

```sql+postgres
select
  domain_name,
  metric_name,
  timestamp,
  average
from
  aws_opensearch_cluster_metric_hourly
where
  average > 100
  AND metric_name = 'CPUUtilization'
order by
  domain_name,
  timestamp;
```

```sql+sqlite
select
  domain_name,
  metric_name,
  timestamp,
  average
from
  aws_opensearch_cluster_metric_hourly
where
  average > Your_Threshold_Value
  AND metric_name = 'Your_Metric_Name'
order by
  domain_name,
  timestamp;
```