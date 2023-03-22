# Table: aws_cloudwatch_metric_statistic_data_point

This table provides statistic data points for the specified metric.

The maximum number of data points returned from a single call is 1,440. If you request more than 1,440 data points, CloudWatch returns an error. To reduce the number of data points, you can narrow the specified time range and make multiple requests across adjacent time ranges, or you can increase the specified period. Data points are not returned in chronological order.

If you need to fetch more than 1440 data points then please use the `aws_cloudwatch_metric_data_point` table.

- You **_must_** specify `metric_name`, and `namespace` in a `where` clause in order to use this table.

- To fetch aggregate statistics data, `dimensions` is not required. However, except for aggregate statistics, you must always pass `dimensions` in the query; the examples below can guide you.

- By default, this table will provide data for the last 24hrs. You can give the `timestamp` value in the below ways to fetch data in a range. The examples below can guide you.

  - timestamp >= ‘2023-03-11T00:00:00Z’ and timestamp <= ‘2023-03-15T00:00:00Z’
  - timestamp between ‘2023-03-11T00:00:00Z’ and ‘2023-03-15T00:00:00Z’
  - timestamp > ‘2023-03-15T00:00:00Z’ (The data will be fetched from the provided time to the current time)
  - timestamp < ‘2023-03-15T00:00:00Z’ (The data will be fetched from one day before the provided time to the provided time)

- We recommend specifying the `period` column in the query to optimize the table output. If you do not specify the `timestamp` then the default value for `period` is 60 seconds. If you specify the `timestamp` then the period will be calculated based on the duration to provide a good spread under the 1440 datapoints.

## Examples

### Aggregate CPU utilization of all EC2 instances for the last 24 hrs

```sql
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

### CPU average utilization of an EC2 instance over 80% for the last 5 days

```sql
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

### Intervals where a volume exceed 1000 average read ops

```sql
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

### CacheHit sum below 10 of an elasticache cluster for for the last 7 days

```sql
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

### Lambda function daily maximum duration over 100 milliseconds

```sql
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

### CPU average utilization of an RDS DB instance over 80% for the last 30 days

```sql
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

### Maximum Bucket size daily statistics of an S3 bucket

```sql
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
