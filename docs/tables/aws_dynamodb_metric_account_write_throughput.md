# Table: aws_dynamodb_metric_account_write_throughput

Amazon CloudWatch metrics provide data about the performance of your systems. The `aws_dynamodb_metric_account_write_throughput` table provides metric statistics at 5 minute intervals for the most recent 5 days.

## Examples

### Basic info

```sql
select
  account_id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_dynamodb_metric_account_write_throughput
order by
  timestamp;
```

### Intervals where throughput exceeds 80 percent

```sql
select
  account_id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_dynamodb_metric_account_write_throughput
where
  maximum > 80
order by
  timestamp;
```
