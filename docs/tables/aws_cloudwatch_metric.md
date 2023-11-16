---
title: "Table: aws_cloudwatch_metric - Query AWS CloudWatch Metrics using SQL"
description: "Allows users to query AWS CloudWatch Metrics to gather information about the performance of their AWS resources and applications."
---

# Table: aws_cloudwatch_metric - Query AWS CloudWatch Metrics using SQL

The `aws_cloudwatch_metric` table in Steampipe provides information about CloudWatch Metrics within AWS CloudWatch. This table allows DevOps engineers to query metric-specific details, including metric names, namespaces, dimensions, and statistics. Users can utilize this table to gather insights on metrics, such as tracking the CPU usage of an EC2 instance, monitoring the latency of an ELB, or even the request count of an API Gateway. The schema outlines the various attributes of the CloudWatch Metric, including the metric name, namespace, dimensions, statistics, and associated metadata.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cloudwatch_metric` table, you can use the `.inspect aws_cloudwatch_metric` command in Steampipe.

**Key columns**:

- `metric_name`: The name of the metric. This is a key identifier and can be used to join with other tables that require metric information.
- `namespace`: The namespace of the metric. This provides context for the metric and can be used to join with tables that require namespace information.
- `dimensions`: The dimensions for the metric. Dimensions are name-value pairs that uniquely identify a metric. They are useful for filtering and aggregating data across various AWS resources.

## Examples

### Basic info

```sql
select
  metric_name,
  namespace,
  dimensions
from
  aws_cloudwatch_metric;
```

### List EBS metrics

```sql
select
  metric_name,
  namespace,
  dimensions
from
  aws_cloudwatch_metric
where
  namespace = 'AWS/EBS';
```

### List EBS `VolumeReadOps` metrics

```sql
select
  metric_name,
  namespace,
  dimensions
from
  aws_cloudwatch_metric
where
  namespace = 'AWS/EBS'
  and metric_name = 'VolumeReadOps';
```

### List metrics for a specific Redshift cluster

```sql
select
  metric_name,
  namespace,
  dimensions
from
  aws_cloudwatch_metric
where
  dimensions_filter = '[
    {"Name": "ClusterIdentifier", "Value": "my-cluster-1"}
  ]'::jsonb;
```

### List EC2 API metrics

```sql
select
  metric_name,
  namespace,
  dimensions
from
  aws_cloudwatch_metric
where
  dimensions_filter = '[
    {"Name": "Type", "Value": "API"},
    {"Name": "Service", "Value": "EC2"}
  ]'::jsonb;
```
