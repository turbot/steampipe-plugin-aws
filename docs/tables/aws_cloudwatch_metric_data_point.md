---
title: "Steampipe Table: aws_cloudwatch_metric_data_point - Query AWS CloudWatch MetricDataPoints using SQL"
description: "Allows users to query AWS CloudWatch MetricDataPoints to fetch detailed information about the data points for a defined metric."
folder: "CloudWatch"
---

# Table: aws_cloudwatch_metric_data_point - Query AWS CloudWatch MetricDataPoints using SQL

The AWS CloudWatch MetricDataPoints is a feature of Amazon CloudWatch that allows you to monitor, store, and access your AWS resources' data in the form of logs and metrics. This feature provides real-time data and insights that can help you optimize the performance and resource utilization of your applications. It also allows you to set alarms and react to changes in your AWS resources, making it easier to troubleshoot issues and discover trends.

## Table Usage Guide

The `aws_cloudwatch_metric_data_point` table in Steampipe provides you with information about MetricDataPoints within AWS CloudWatch. This table enables you, as a DevOps engineer, to query specific details about the data points for a defined metric, including the timestamp, sample count, sum, minimum, and maximum values. You can utilize this table to gather insights on metrics, such as tracking the number of requests to an application over time, monitoring the CPU usage and network traffic of EC2 instances, and more. The schema outlines the various attributes of the MetricDataPoint, including the average, sample count, sum, minimum, and maximum values, along with the timestamp of the data point.

**Important Notes**
This table provides metric data points for the specified id. The maximum number of data points returned from a single call is 100,800.

- You must specify `id`, and `expression` or `id`, and `metric_stat` in a `where` clause in order to use this table.
- By default, this table will provide data for the last 24hrs. You can provide the `timestamp` value in the following ways to fetch data in a range. The examples below can guide you.
  - timestamp >= ‘2023-03-11T00:00:00Z’ and timestamp <= ‘2023-03-15T00:00:00Z’
  - timestamp between ‘2023-03-11T00:00:00Z’ and ‘2023-03-15T00:00:00Z’
  - timestamp > ‘2023-03-15T00:00:00Z’ (The data will be fetched from the provided time to the current time)
  - timestamp < ‘2023-03-15T00:00:00Z’ (The data will be fetched from one day before the provided time to the provided time)
- It's recommended that you specify the `period` column in the query to optimize the table output. If you do not specify the `timestamp` then the default value for `period` is 60 seconds. If you specify the `timestamp` then the period will be calculated based on the duration mentioned ([here](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/cloudwatch/types#MetricStat.Period)).
- Using this table adds to the cost of your monthly bill from AWS. Optimizations have

## Examples

### Aggregate maximum CPU utilization of all EC2 instances for the last 24 hrs
Determine the peak CPU usage of all EC2 instances in the past day. This query can be used to monitor system performance and identify potential issues related to high CPU utilization.

```sql+postgres
select
  id,
  label,
  timestamp,
  period,
  value,
  expression
from
  aws_cloudwatch_metric_data_point
where
  id = 'm1'
  and expression = 'select max(CPUUtilization) from schema("AWS/EC2", InstanceId)'
order by
  timestamp;
```

```sql+sqlite
select
  id,
  label,
  timestamp,
  period,
  value,
  expression
from
  aws_cloudwatch_metric_data_point
where
  id = 'm1'
  and expression = 'select max(CPUUtilization) from schema("AWS/EC2", InstanceId)'
order by
  timestamp;
```

### Calculate error rate on the provided custom metric ID for the last 24 hrs
This query is useful for monitoring the error rate on a specific custom metric over the last 24 hours. It can help identify potential issues or anomalies in your system, allowing for timely troubleshooting and maintenance.

```sql+postgres
select
  id,
  label,
  timestamp,
  period,
  value,
  expression
from
  aws_cloudwatch_metric_data_point
where
  id = 'e1'
  and expression = 'SUM(METRICS(''error''))'
order by
  timestamp;
```

```sql+sqlite
select
  id,
  label,
  timestamp,
  period,
  value,
  expression
from
  aws_cloudwatch_metric_data_point
where
  id = 'e1'
  and expression = 'SUM(METRICS(''error''))'
order by
  timestamp;
```

### CPU average utilization of multiple EC2 instances over 80% for the last 5 days
Identify instances where the average CPU utilization of multiple EC2 instances has exceeded 80% in the past 5 days. This can be useful in monitoring resource usage and identifying potential performance issues.

```sql+postgres
select
  id,
  label,
  timestamp,
  period,
  round(value::numeric, 2) as avg_cpu,
  metric_stat
from
  aws_cloudwatch_metric_data_point
where
  id = 'm1'
  and value > 80
  and timestamp >= now() - interval '5 day'
  and metric_stat = '{
    "Metric": {
    "Namespace": "AWS/EC2",
    "MetricName": "CPUUtilization",
    "Dimensions": [
      {
        "Name": "InstanceId",
        "Value": "i-0353536c53f7c8235"
      },
      {
        "Name": "InstanceId",
        "Value": "i-0dd7043e0f6f0f36d"
      }
    ]},
    "Stat": "Average"}'
order by
  timestamp;
```

```sql+sqlite
select
  id,
  label,
  timestamp,
  period,
  round(value, 2) as avg_cpu,
  metric_stat
from
  aws_cloudwatch_metric_data_point
where
  id = 'm1'
  and value > 80
  and timestamp >= datetime('now','-5 day')
  and json_extract(metric_stat, '$.Metric.Namespace') = 'AWS/EC2'
  and json_extract(metric_stat, '$.Metric.MetricName') = 'CPUUtilization'
  and json_extract(metric_stat, '$.Metric.Dimensions[0].Name') = 'InstanceId'
  and json_extract(metric_stat, '$.Metric.Dimensions[0].Value') = 'i-0353536c53f7c8235'
  and json_extract(metric_stat, '$.Metric.Dimensions[1].Name') = 'InstanceId'
  and json_extract(metric_stat, '$.Metric.Dimensions[1].Value')
```

### Intervals where an EBS volume exceed 1000 average read ops daily
Explore instances where an EBS volume exceeds a daily average of 1000 read operations. This can be useful in understanding the performance and load on your EBS volumes, helping you make informed decisions about capacity planning and resource allocation.

```sql+postgres
select
  id,
  label,
  timestamp,
  value,
  metric_stat
from
  aws_cloudwatch_metric_data_point
where
  id = 'm1'
  and value > 1000
  and period = 86400
  and scan_by = 'TimestampDescending'
  and timestamp between '2023-03-10T00:00:00Z' and '2023-03-16T00:00:00Z'
  and metric_stat = '{
    "Metric": {
    "Namespace": "AWS/EBS",
    "MetricName": "VolumeReadOps",
    "Dimensions": [
      {
        "Name": "VolumeId",
        "Value": "vol-00607053b218c6d74"
      }
    ]},
    "Stat": "Average"}';
```

```sql+sqlite
select
  id,
  label,
  timestamp,
  value,
  metric_stat
from
  aws_cloudwatch_metric_data_point
where
  id = 'm1'
  and value > 1000
  and period = 86400
  and scan_by = 'TimestampDescending'
  and timestamp between '2023-03-10T00:00:00Z' and '2023-03-16T00:00:00Z'
  and json_extract(metric_stat, '$.Metric.Namespace') = 'AWS/EBS'
  and json_extract(metric_stat, '$.Metric.MetricName') = 'VolumeReadOps'
  and json_extract(metric_stat, '$.Metric.Dimensions[0].Name') = 'VolumeId'
  and json_extract(metric_stat, '$.Metric.Dimensions[0].Value') = 'vol-00607053b218c6d74'
  and json_extract(metric_stat, '$.Stat
```

### CacheHit sum below 10 of an elasticache cluster for the last 7 days
Determine the performance of an ElastiCache cluster over the past week by identifying instances where cache hit sums were less than 10. This can be useful for analyzing the effectiveness of your cache configuration and identifying potential areas for improvement.

```sql+postgres
select
  id,
  label,
  timestamp,
  value,
  metric_stat
from
  aws_cloudwatch_metric_data_point
where
  id = 'e1'
  and value < 10
  and timestamp >= now() - interval '7 day'
  and metric_stat = '{
    "Metric": {
    "Namespace": "AWS/ElastiCache",
    "MetricName": "CacheHits",
    "Dimensions": [
      {
        "Name": "CacheClusterId",
        "Value": "cluster-delete-001"
      }
    ]},
    "Stat": "Sum"}'
order by
  timestamp;
```

```sql+sqlite
select
  id,
  label,
  timestamp,
  value,
  metric_stat
from
  aws_cloudwatch_metric_data_point
where
  id = 'e1'
  and value < 10
  and timestamp >= datetime('now', '-7 days')
  and metric_stat = '{
    "Metric": {
    "Namespace": "AWS/ElastiCache",
    "MetricName": "CacheHits",
    "Dimensions": [
      {
        "Name": "CacheClusterId",
        "Value": "cluster-delete-001"
      }
    ]},
    "Stat": "Sum"}'
order by
  timestamp;
```

### Maximum Bucket size daily statistics of an S3 bucket for an account
Explore the maximum storage usage of a specific S3 bucket in your AWS account within a specific timeframe. This can help manage storage capacity and understand usage patterns.

```sql+postgres
select
  id,
  label,
  timestamp,
  value,
  metric_stat
from
  aws_cloudwatch_metric_data_point
where
  id = 'e1'
  and source_account_id = '533743456432100'
  and timestamp between '2023-03-10T00:00:00Z' and '2023-03-16T00:00:00Z'
  and metric_stat = '{
    "Metric": {
    "Namespace": "AWS/S3",
    "MetricName": "BucketSizeBytes",
    "Dimensions": [
      {
        "Name": "BucketName",
        "Value": "steampipe-test"
      },
      {
        "Name": "StorageType",
        "Value": "StandardStorage"
      }
    ]},
    "Stat": "Maximum"}'
order by
  timestamp;
```

```sql+sqlite
select
  id,
  label,
  timestamp,
  value,
  metric_stat
from
  aws_cloudwatch_metric_data_point
where
  id = 'e1'
  and source_account_id = '533743456432100'
  and timestamp between '2023-03-10T00:00:00Z' and '2023-03-16T00:00:00Z'
  and json_extract(metric_stat, '$.Metric.Namespace') = 'AWS/S3'
  and json_extract(metric_stat, '$.Metric.MetricName') = 'BucketSizeBytes'
  and json_extract(metric_stat, '$.Metric.Dimensions[0].Name') = 'BucketName'
  and json_extract(metric_stat, '$.Metric.Dimensions[0].Value') = 'steampipe-test'
  and json_extract(metric_stat, '$.Metric.Dimensions[1].Name') = 'StorageType'
  and json_extract(metric_stat
```