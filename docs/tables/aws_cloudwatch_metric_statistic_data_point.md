# Table: aws_cloudwatch_metric_statistic_data_point

This table provides statistic data points for the specified metric.

The maximum number of data points returned from a single call is 1,440. If you request more than 1,440 data points, CloudWatch returns an error. To reduce the number of data points, you can narrow the specified time range and make multiple requests across adjacent time ranges, or you can increase the specified period. Data points are not returned in chronological order.

If you need to fetch more than 1440 data points then please use the `aws_cloudwatch_metric_data_point` table.

- You **_must_** specify `metric_name`, `namespace`, `timestamp` and `dimensions` in a `where` clause in order to use this table. To fetch aggregate statistics data `dimensions` is not required.

- You need to provide the timestamp value in a range like `timestamp >= '2023-03-15T00:00:00Z' and timestamp <= '2023-03-15T13:33:00Z' or like in the examples below.`

- The GetMetricStatistics API used for this table cannot process multiple dimension values at a time, so you need to pass one set of dimensions to the query, like in the examples below.

We recommend specifying the `period` column in the query to optimize the table output.

## Examples

### Aggregate CPU utilization of all ec2 instances for a given time frame

```sql
select
  metric_name,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_cloudwatch_metric_statistic_data_point
where
  namespace = 'AWS/EC2'
  and metric_name = 'CPUUtilization'
  and timestamp between '2023-03-15T00:00:00Z' and '2023-03-15T13:33:00Z'
  and period = 86400;
```

### CPU average utilization of an ec2 instance over 80% for a given time frame

```sql
select
  jsonb_pretty(dimensions) as dimensions,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_cloudwatch_metric_statistic_data_point
where
  namespace = 'AWS/EC2'
  and metric_name = 'CPUUtilization'
  and average > 80
  and timestamp between '2023-03-15T00:00:00Z' and '2023-03-15T13:33:00Z'
  and period = 86400
  and dimensions = '[
    {"Name": "InstanceId", "Value": "i-0dd7043e0f6f0f36d"}
    ]'::jsonb
order by
  timestamp;
```

### Intervals where a volume exceed 1000 average read ops for a given time frame

```sql
select
  jsonb_pretty(dimensions) as dimensions,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_cloudwatch_metric_statistic_data_point
where
  namespace = 'AWS/EBS'
  and metric_name = 'VolumeReadOps'
  and average > 1000
  and timestamp between '2023-03-15T00:00:00Z' and '2023-03-15T13:33:00Z'
  and period = 86400
  and dimensions = '[
    {"Name": "VolumeId", "Value": "vol-00607053b218c6d74"}
    ]'::jsonb
order by
  timestamp;
```

### CacheHit sum below 10 of an elasticache cluster for a given time frame

```sql
select
  jsonb_pretty(dimensions) as dimensions,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_cloudwatch_metric_statistic_data_point
where
  namespace = 'AWS/ElastiCache'
  and metric_name = 'CacheHits'
  and sum < 10
  and timestamp between '2023-03-15T00:00:00Z' and '2023-03-15T13:33:00Z'
  and period = 86400
  and dimensions = '[
    {"Name": "CacheClusterId", "Value": "cluster-delete-001"}
    ]'::jsonb
order by
  timestamp;
```

### Lambda function daily maximum duration over 100 milliseconds for a given time frame

```sql
select
  jsonb_pretty(dimensions) as dimensions,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_cloudwatch_metric_statistic_data_point
where
  namespace = 'AWS/Lambda'
  and metric_name = 'Duration'
  and maximum > 100
  and timestamp between '2023-03-15T00:00:00Z' and '2023-03-15T13:33:00Z'
  and period = 86400
  and dimensions = '[
    {"Name": "FunctionName", "Value": "test"}
    ]'::jsonb
order by
  timestamp;
```

### CPU average utilization of an RDS DB instance over 80% for a given time frame

```sql
select
  jsonb_pretty(dimensions) as dimensions,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_cloudwatch_metric_statistic_data_point
where
  namespace = 'AWS/RDS'
  and metric_name = 'CPUUtilization'
  and average > 80
  and timestamp between '2023-03-15T00:00:00Z' and '2023-03-15T13:33:00Z'
  and period = 86400
  and dimensions = '[
    {"Name": "DBInstanceIdentifier", "Value": "database-1"}
    ]'::jsonb
order by
  timestamp;
```

### Bucket size statistics of an s3 bucket for a given time frame

```sql
select
  jsonb_pretty(dimensions) as dimensions,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_cloudwatch_metric_statistic_data_point
where
  namespace = 'AWS/S3'
  and metric_name = 'BucketSizeBytes'
  and timestamp between '2023-03-6T00:00:00Z' and '2023-03-15T13:33:00Z'
  and period = 86400
  and dimensions = '[
    {"Name": "BucketName", "Value": "steampipe-test"},
    {"Name": "StorageType", "Value": "StandardStorage"}
    ]'::jsonb
order by
  timestamp;
```
