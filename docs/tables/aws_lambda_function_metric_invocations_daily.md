# Table: aws_lambda_function_metric_invocations_daily

Amazon CloudWatch Metrics provide data about the performance of your systems.  The `aws_lambda_function_metric_invocations_daily` table provides metric statistics at 24 hour intervals for the last year.

## Examples

### Basic info

```sql
select
  name,
  timestamp,
  sum
from
  aws_lambda_function_metric_invocations_daily
order by
  name,
  timestamp;
```



### Lambda function daily invocations over 10 in last 3 days

```sql
select
  name,
  timestamp,
  round(sum::numeric,2) as sum_invocations,
  sample_count
from
  aws_lambda_function_metric_invocations_daily
where 
    date_part('day', now() - timestamp) <=3
and sum > 10
order by
  name,
  timestamp;
```