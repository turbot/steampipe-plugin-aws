# Table: aws_ecs_cluster_metric_cpu_utilization_hourly

Amazon CloudWatch metrics provide data about the performance of your systems. The `aws_ecs_cluster_metric_cpu_utilization_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

## Examples

### Basic info

```sql
select
  cluster_name,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_ecs_cluster_metric_cpu_utilization_hourly
order by
  cluster_name,
  timestamp;
```

### CPU Over 80% average

```sql
select
  cluster_name,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_ecs_cluster_metric_cpu_utilization_hourly
where
  average > 80
order by
  cluster_name,
  timestamp;
```

### CPU hourly average < 1%

```sql
select
  cluster_name,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_ecs_cluster_metric_cpu_utilization_hourly
where
  average < 1
order by
  cluster_name,
  timestamp;
```
