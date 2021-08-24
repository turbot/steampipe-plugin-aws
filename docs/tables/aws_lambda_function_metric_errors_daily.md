# Table: aws_lambda_function_metric_errors_daily

Amazon CloudWatch Metrics provide data about the performance of your systems. The `aws_lambda_function_metric_errors_daily` table provides metric statistics at 24 hour intervals for the last year.

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
  aws_lambda_function_metric_errors_daily
order by
  name,
  timestamp;
```

### Lambda function daily average error less than 1

```sql
select
  name,
  timestamp,
  round(minimum::numeric,2) as min_error,
  round(maximum::numeric,2) as max_error,
  round(average::numeric,2) as avg_error,
  sample_count
from
  aws_lambda_function_metric_errors_daily
where average < 1
order by
  name,
  timestamp;
```