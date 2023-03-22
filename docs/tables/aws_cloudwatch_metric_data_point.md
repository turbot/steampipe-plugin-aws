# Table: aws_cloudwatch_metric_data_point

This table provides metric data points for the specified id. The maximum number of data points returned from a single call is 100,800.

- You **_must_** specify `id`, and `expression` or `id`, and `metric_stat` in a `where` clause in order to use this table.

- By default, this table will provide data for the last 24hrs. You can give the `timestamp` value in the below ways to fetch data in a range. The examples below can guide you.

  - timestamp >= ‘2023-03-11T00:00:00Z’ and timestamp <= ‘2023-03-15T00:00:00Z’
  - timestamp between ‘2023-03-11T00:00:00Z’ and ‘2023-03-15T00:00:00Z’
  - timestamp > ‘2023-03-15T00:00:00Z’ (The data will be fetched from the provided time to the current time)
  - timestamp < ‘2023-03-15T00:00:00Z’ (The data will be fetched from one day before the provided time to the provided time)

- We recommend specifying the `period` column in the query to optimize the table output. If you do not specify the `timestamp` then the default value for `period` is 60 seconds. If you specify the `timestamp` then the period will be calculated based on the duration mentioned ([here](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/cloudwatch/types#MetricStat.Period)).

Note: Using this table adds to cost to your monthly bill from AWS. Optimizations have been put in place to minimize the impact as much as possible. Please refer to AWS CloudWatch Pricing to understand the cost implications.

## Examples

### Aggregate maximum CPU utilization of all EC2 instances for the last 24 hrs

```sql
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

```sql
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

```sql
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

### Intervals where a volume exceed 1000 average read ops daily

```sql
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

### CacheHit sum below 10 of an elasticache cluster for the last 7 days

```sql
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

### Maximum Bucket size daily statistics of an S3 bucket for an account

```sql
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
  and account_id = '533743456432100'
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
