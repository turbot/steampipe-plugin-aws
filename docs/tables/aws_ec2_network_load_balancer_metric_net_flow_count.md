# Table: aws_ec2_network_load_balancer_metric_net_flow_count

Amazon CloudWatch Metrics provide data about the performance of your systems. The `aws_ec2_network_load_balancer_metric_net_flow_count` table provides metric statistics at 5 min intervals for the most recent 5 days.

## Examples

### Basic info

```sql
select
  name,
  metric_name,
  namespace,
  maximum,
  minimum,
  sample_count,
  timestamp
from
  aws_ec2_network_load_balancer_metric_net_flow_count
order by
  name,
  timestamp;
```

### Intervals where net flow count < 100

```sql
select
  name,
  metric_name,
  namespace,
  maximum,
  minimum,
  average,
  sample_count,
  timestamp
from
  aws_ec2_network_load_balancer_metric_net_flow_count
where
  average < 100
order by
  name,
  timestamp;
```
