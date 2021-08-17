# Table: aws_ec2_application_load_balancer_metric_request_count_daily

Amazon CloudWatch Metrics provide data about the performance of your systems. The `aws_ec2_application_load_balancer_metric_request_count_daily` table provides metric statistics at 24 hour intervals for the most recent 1 year.

## Examples

### Basic info

```sql
select
  name,
  metric_name,
  namespace,
  average,
  maximum,
  minimum,
  sample_count,
  timestamp
from
  aws_ec2_application_load_balancer_metric_request_count_daily
order by
  name,
  timestamp;
```

### Intervals averaging less than 100 request count

```sql
select
  name,
  metric_name,
  namespace,
  maximum,
  minimum,
  average
  sample_count,
  timestamp
from
  aws_ec2_application_load_balancer_metric_request_count_daily
where
  average < 100
order by
  name,
  timestamp;
```
