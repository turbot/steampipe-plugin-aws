---
title: "Steampipe Table: aws_cloudwatch_metric_statistic_data_point - Query AWS CloudWatch Metric Statistics Data Point using SQL"
description: "Allows users to query AWS CloudWatch Metric Statistics Data Point to obtain detailed metrics data."
folder: "CloudWatch"
---

# Table: aws_cloudwatch_metric_statistic_data_point - Query AWS CloudWatch Metric Statistics Data Point using SQL

The AWS CloudWatch Metric Statistics Data Point is a feature of the Amazon CloudWatch service. It allows you to retrieve statistical data about your AWS resources that is collected by CloudWatch. This statistical data can be used for monitoring, troubleshooting, and setting alarms for when specific thresholds are met.

## Table Usage Guide

The `aws_cloudwatch_metric_statistic_data_point` table in Steampipe provides you with information about the data points for a specified metric in AWS CloudWatch. This table allows you, as a DevOps engineer, to query detailed metrics data, including timestamps, samples count, maximum, minimum, and average values. You can utilize this table to gather insights on metric data points, such as observing trends, identifying peaks or anomalies, and monitoring the overall performance of AWS resources. The schema outlines the various attributes of the metric data points, including the namespace, metric name, dimensions, and the period, start time, and end time of the data points.

**Important Notes**
The maximum number of data points that can be returned from a single call is 1,440. If you request more than 1,440 data points, CloudWatch will return an error. To reduce the number of data points, you can narrow the specified time range and make multiple requests across adjacent time ranges, or you can increase the specified period. Please note that data points are not returned in chronological order.

- If you need to fetch more than 1440 data points, you should use the `aws_cloudwatch_metric_data_point` table.

- You must specify `metric_name`, and `namespace` in a `where` clause in order to use this table.

- To fetch aggregate statistics data, `dimensions` is not required. However, except for aggregate statistics, you must always pass `dimensions` in the query; the examples below can guide you.

- By default, this table will provide data for the last 24hrs. You can give the `timestamp` value in the below ways to fetch data in a range. The examples below can guide you.

  - timestamp >= ‘2023-03-11T00:00:00Z’ and timestamp <= ‘2023-03-15T00:00:00Z’
  - timestamp between ‘2023-03-11T00:00:00Z’ and ‘2023-03-15T00:00:00Z’
  - timestamp > ‘2023-03-15T00:00:00Z’ (The data will be fetched from the provided time to the current time)
  - timestamp < ‘2023-03-15T00:00:00Z’ (The data will be fetched from one day before the provided time to the provided time)

- We recommend specifying the `period` column in the query to optimize the table output. If you do not specify the `timestamp` then the default value for `period` is 60 seconds. If you specify the `timestamp` then the period will be calculated based on the duration to provide a good spread under the 1440 datapoints.

## Examples

### Aggregate CPU utilization of all EC2 instances for the last 24 hrs
Explore the extent of CPU usage across all EC2 instances over the past day. This is useful to monitor system performance, identify potential bottlenecks, and plan for capacity upgrades.

```sql+postgres
select
  metric_name,
  timestamp,
  round(minimum::numeric, 2) as min_cpu,
  round(maximum::numeric, 2) as max_cpu,
  round(average::numeric, 2) as avg_cpu,
  sum,
  sample_count
from
  aws_cloudwatch_metric_statistic_data_point
where
  namespace = 'AWS/EC2'
  and metric_name = 'CPUUtilization'
order by
  timestamp;
```

```sql+sqlite
select
  metric_name,
  timestamp,
  round(minimum, 2) as min_cpu,
  round(maximum, 2) as max_cpu,
  round(average, 2) as avg_cpu,
  sum,
  sample_count
from
  aws_cloudwatch_metric_statistic_data_point
where
  namespace = 'AWS/EC2'
  and metric_name = 'CPUUtilization'
order by
  timestamp;
```

### CPU average utilization of an EC2 instance over 80% for the last 5 days
Determine the instances where the average utilization of a specific EC2 instance has exceeded 80% in the past 5 days. This query is useful for monitoring system performance and identifying potential issues with resource allocation.

```sql+postgres
select
  jsonb_pretty(dimensions) as dimensions,
  timestamp,
  round(average::numeric, 2) as avg_cpu
from
  aws_cloudwatch_metric_statistic_data_point
where
  namespace = 'AWS/EC2'
  and metric_name = 'CPUUtilization'
  and average > 80
  and timestamp >= now() - interval '5 day'
  and dimensions = '[
    {"Name": "InstanceId", "Value": "i-0dd7043e0f6f0f36d"}
    ]'
order by
  timestamp;
```

```sql+sqlite
select
  json_pretty(dimensions) as dimensions,
  timestamp,
  round(average, 2) as avg_cpu
from
  aws_cloudwatch_metric_statistic_data_point
where
  namespace = 'AWS/EC2'
  and metric_name = 'CPUUtilization'
  and average > 80
  and timestamp >= datetime('now', '-5 day')
  and dimensions = '[
    {"Name": "InstanceId", "Value": "i-0dd7043e0f6f0f36d"}
    ]'
order by
  timestamp;
```

### Intervals where a volume exceed 1000 average read ops
Identify instances where the average read operations on a specific volume surpasses a set threshold within a defined timeframe. This query aids in analyzing periods of high read operations, helping to optimize resource usage and performance in AWS EBS.

```sql+postgres
select
  jsonb_pretty(dimensions) as dimensions,
  timestamp,
  average
from
  aws_cloudwatch_metric_statistic_data_point
where
  namespace = 'AWS/EBS'
  and metric_name = 'VolumeReadOps'
  and average > 1000
  and timestamp between '2023-03-10T00:00:00Z' and '2023-03-16T00:00:00Z'
  and period = 300
  and dimensions = '[
    {"Name": "VolumeId", "Value": "vol-00607053b218c6d74"}
    ]'
order by
  timestamp;
```

```sql+sqlite
select
  dimensions,
  timestamp,
  average
from
  aws_cloudwatch_metric_statistic_data_point
where
  namespace = 'AWS/EBS'
  and metric_name = 'VolumeReadOps'
  and average > 1000
  and timestamp between '2023-03-10T00:00:00Z' and '2023-03-16T00:00:00Z'
  and period = 300
  and json_extract(dimensions, '

### CacheHit sum below 10 of an elasticache cluster for for the last 7 days
Analyze the performance of an ElastiCache cluster by tracking instances where cache hits were less than 10 over the past week. This could be useful to identify potential issues with the cache's configuration or usage patterns.

```sql+postgres
select
  jsonb_pretty(dimensions) as dimensions,
  timestamp,
  sum
from
  aws_cloudwatch_metric_statistic_data_point
where
  namespace = 'AWS/ElastiCache'
  and metric_name = 'CacheHits'
  and sum < 10
  and timestamp >= now() - interval '7 day'
  and dimensions = '[
    {"Name": "CacheClusterId", "Value": "cluster-delete-001"}
    ]'
order by
  timestamp;
```

```sql+sqlite
select
  json_pretty(dimensions) as dimensions,
  timestamp,
  sum
from
  aws_cloudwatch_metric_statistic_data_point
where
  namespace = 'AWS/ElastiCache'
  and metric_name = 'CacheHits'
  and sum < 10
  and timestamp >= datetime('now', '-7 day')
  and dimensions = '[
    {"Name": "CacheClusterId", "Value": "cluster-delete-001"}
    ]'
order by
  timestamp;
```

### Lambda function daily maximum duration over 100 milliseconds
Discover the instances when your AWS Lambda function's maximum daily duration exceeds 100 milliseconds within a specific time frame. This can be useful for identifying potential performance issues or bottlenecks in your application.

```sql+postgres
select
  jsonb_pretty(dimensions) as dimensions,
  timestamp,
  maximum
from
  aws_cloudwatch_metric_statistic_data_point
where
  namespace = 'AWS/Lambda'
  and metric_name = 'Duration'
  and maximum > 100
  and timestamp >= '2023-02-15T00:00:00Z'
  and timestamp <= '2023-03-15T00:00:00Z'
  and period = 86400
  and dimensions = '[
    {"Name": "FunctionName", "Value": "test"}
    ]'
order by
  timestamp;
```

```sql+sqlite
select
  json_pretty(dimensions) as dimensions,
  timestamp,
  maximum
from
  aws_cloudwatch_metric_statistic_data_point
where
  namespace = 'AWS/Lambda'
  and metric_name = 'Duration'
  and maximum > 100
  and timestamp >= '2023-02-15T00:00:00Z'
  and timestamp <= '2023-03-15T00:00:00Z'
  and period = 86400
  and dimensions = '[
    {"Name": "FunctionName", "Value": "test"}
    ]'
order by
  timestamp;
```

### CPU average utilization of an RDS DB instance over 80% for the last 30 days
This query is used to monitor the performance of an RDS DB instance by tracking its CPU utilization. If the average CPU usage exceeds 80% over the past 30 days, it may indicate a need for more resources or optimization to prevent potential system slowdowns or failures.

```sql+postgres
select
  jsonb_pretty(dimensions) as dimensions,
  timestamp,
  round(average::numeric, 2) as avg_cpu
from
  aws_cloudwatch_metric_statistic_data_point
where
  namespace = 'AWS/RDS'
  and metric_name = 'CPUUtilization'
  and average > 80
  and timestamp >= now() - interval '30 day'
  and dimensions = '[
    {"Name": "DBInstanceIdentifier", "Value": "database-1"}
    ]'
order by
  timestamp;
```

```sql+sqlite
select
  json_pretty(dimensions) as dimensions,
  timestamp,
  round(average, 2) as avg_cpu
from
  aws_cloudwatch_metric_statistic_data_point
where
  namespace = 'AWS/RDS'
  and metric_name = 'CPUUtilization'
  and average > 80
  and timestamp >= datetime('now', '-30 day')
  and dimensions = '[
    {"Name": "DBInstanceIdentifier", "Value": "database-1"}
    ]'
order by
  timestamp;
```

### Maximum Bucket size daily statistics of an S3 bucket
Explore the daily storage usage of a specific S3 bucket over a given time period. This is useful for tracking storage trends and planning for future capacity needs.

```sql+postgres
select
  jsonb_pretty(dimensions) as dimensions,
  timestamp,
  minimum
from
  aws_cloudwatch_metric_statistic_data_point
where
  namespace = 'AWS/S3'
  and metric_name = 'BucketSizeBytes'
  and timestamp between '2023-03-6T00:00:00Z' and '2023-03-15T00:00:00Z'
  and period = 86400
  and dimensions = '[
    {"Name": "BucketName", "Value": "steampipe-test"},
    {"Name": "StorageType", "Value": "StandardStorage"}
    ]'
order by
  timestamp;
```

```sql+sqlite
select
  json_pretty(dimensions) as dimensions,
  timestamp,
  minimum
from
  aws_cloudwatch_metric_statistic_data_point
where
  namespace = 'AWS/S3'
  and metric_name = 'BucketSizeBytes'
  and timestamp between '2023-03-6T00:00:00Z' and '2023-03-15T00:00:00Z'
  and period = 86400
  and dimensions = '[
    {"Name": "BucketName", "Value": "steampipe-test"},
    {"Name": "StorageType", "Value": "StandardStorage"}
    ]'
order by
  timestamp;
```
