# Table: aws_lambda_function_metric_duration_daily

Amazon CloudWatch Metrics provide data about the performance of your systems.  The `aws_lambda_function_metric_duration_daily` table provides metric statistics at 24 hour intervals for the last year.


## Examples

### Basic info

```sql
select
  name,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_lambda_function_metric_duration_daily
order by
  name,
  timestamp;
```

### Lambda function daily maximum duration over 100 milliseconds

```sql
select
  name,
  timestamp,
  round(minimum::numeric,2) as min_duration,
  round(maximum::numeric,2) as max_duration,
  round(average::numeric,2) as avg_duration,
  sample_count
from
  aws_lambda_function_metric_duration_daily
where maximum > 100
order by
  name,
  timestamp;
```

### Lambda function daily average duration less than 5 milliseconds

```sql
select
  name,
  timestamp,
  round(minimum::numeric,2) as min_duration,
  round(maximum::numeric,2) as max_duration,
  round(average::numeric,2) as avg_duration,
  sample_count
from
  aws_lambda_function_metric_duration_daily
where average < 5
order by
  name,
  timestamp;
```