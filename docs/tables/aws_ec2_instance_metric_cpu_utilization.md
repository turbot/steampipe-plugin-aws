# Table: aws_ec2_instance_metric_cpu_utilization

Amazon CloudWatch Metrics provide data about the performance of your systems.  The `aws_ec2_instance_metric_cpu_utilization` table provides metric statistics at 5 minute intervals for the most recent 5 days.


## Examples

### Basic info

```sql
select
  instance_id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_ec2_instance_metric_cpu_utilization
order by
  instance_id,
  timestamp;
```



### CPU Over 80% average

```sql
select
  instance_id,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_ec2_instance_metric_cpu_utilization
where average > 80
order by
  instance_id,
  timestamp;
```