---
title: "Steampipe Table: aws_cloudwatch_metric - Query AWS CloudWatch Metrics using SQL"
description: "Allows users to query AWS CloudWatch Metrics to gather information about the performance of their AWS resources and applications."
folder: "CloudWatch"
---

# Table: aws_cloudwatch_metric - Query AWS CloudWatch Metrics using SQL

The AWS CloudWatch Metrics is a feature of Amazon CloudWatch that allows you to monitor, store, and access your log files from Amazon Elastic Compute Cloud (EC2) instances, AWS CloudTrail, Route 53, and other sources. It provides data and actionable insights to monitor your applications, understand and respond to system-wide performance changes, optimize resource utilization, and get a unified view of operational health. By using SQL queries with CloudWatch Metrics, you can gain a deeper understanding of your system's operational status.

## Table Usage Guide

The `aws_cloudwatch_metric` table in Steampipe provides you with information about CloudWatch Metrics within AWS CloudWatch. This table allows you, as a DevOps engineer, to query metric-specific details, including metric names, namespaces, dimensions, and statistics. You can utilize this table to gather insights on metrics, such as tracking the CPU usage of an EC2 instance, monitoring the latency of an ELB, or even the request count of an API Gateway. The schema outlines the various attributes of the CloudWatch Metric for you, including the metric name, namespace, dimensions, statistics, and associated metadata.

**Important Notes**
- You can include up to 10 dimensions in the `dimensions_filter` column.

## Examples

### Basic info
Explore the metrics and their associated namespaces in your AWS CloudWatch service. This can help you understand the different performance indicators being monitored and their corresponding AWS services, providing a comprehensive overview of your system's performance and health.

```sql+postgres
select
  metric_name,
  namespace,
  dimensions
from
  aws_cloudwatch_metric;
```

```sql+sqlite
select
  metric_name,
  namespace,
  dimensions
from
  aws_cloudwatch_metric;
```

### List EBS metrics
Explore the performance metrics related to Amazon Elastic Block Store (EBS) to gain insights into its operations and efficiency. This can help in identifying potential issues and optimizing resource usage.

```sql+postgres
select
  metric_name,
  namespace,
  dimensions
from
  aws_cloudwatch_metric
where
  namespace = 'AWS/EBS';
```

```sql+sqlite
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
Discover the segments that track the read operations on your Elastic Block Store (EBS) volumes. This is useful for monitoring the performance and usage patterns of your EBS volumes in AWS environment.

```sql+postgres
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

```sql+sqlite
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
Explore the performance metrics of a specific Redshift cluster to gain insights into its operational efficiency and resource utilization. This can be useful in monitoring the cluster's health and optimizing its performance.

```sql+postgres
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

```sql+sqlite
select
  metric_name,
  namespace,
  dimensions
from
  aws_cloudwatch_metric
where
  json_extract(dimensions_filter, '$[0].Name') = 'ClusterIdentifier' 
  and json_extract(dimensions_filter, '$[0].Value') = 'my-cluster-1';
```

### List EC2 API metrics
Explore which API metrics are available for the EC2 service in AWS Cloudwatch. This is useful for monitoring and optimizing the performance of your EC2 instances.

```sql+postgres
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

```sql+sqlite
select
  metric_name,
  namespace,
  dimensions
from
  aws_cloudwatch_metric
where
  json_extract(dimensions_filter, '$[0].Name') = "Type" and json_extract(dimensions_filter, '$[0].Value') = "API"
  and json_extract(dimensions_filter, '$[1].Name') = "Service" and json_extract(dimensions_filter, '$[1].Value') = "EC2";
```