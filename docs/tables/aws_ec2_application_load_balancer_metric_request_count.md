# Table: aws_ec2_application_load_balancer_metric_request_count

Amazon CloudWatch Metrics provide data about the performance of your systems.  The `aws_ec2_application_load_balancer_metric_request_count` table provides metric statistics at 5 min intervals for the most recent 7 days.

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
  aws_ec2_application_load_balancer_metric_request_count
order by
  name,
  timestamp;
```

### Intervals averaging less than 100 net flow count

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
  aws_ec2_application_load_balancer_metric_request_count
where
  average < 100
order by
  name,
  timestamp;
```