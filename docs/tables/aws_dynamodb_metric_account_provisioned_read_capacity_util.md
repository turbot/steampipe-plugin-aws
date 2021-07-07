# Table: aws_dynamodb_metric_account_provisioned_read_capacity_util

Amazon CloudWatch Metrics provide data about the performance of your systems. The `aws_dynamodb_metric_account_provisioned_read_capacity_util` table provides metric statistics at 5 minute intervals for the most recent 5 days.

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
  aws_dynamodb_metric_account_provisioned_read_capacity_util
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
  aws_dynamodb_metric_account_provisioned_read_capacity_util
where
  maximum > 80
order by
  timestamp;
```
