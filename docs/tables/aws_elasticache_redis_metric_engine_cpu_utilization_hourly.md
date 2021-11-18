# Table: aws_elasticache_redis_metric_engine_cpu_utilization_hourly

Amazon CloudWatch Metrics provide data about the performance of your systems.The `aws_elasticache_redis_metric_engine_cpu_utilization_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

## Examples

### Basic info

```sql
select
  cacheclusterid,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_elasticache_redis_metric_engine_cpu_utilization_hourly
order by
  cacheclusterid,
  timestamp;
```

### CPU Over 80% average

```sql
select
  cacheclusterid,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_elasticache_redis_metric_engine_cpu_utilization_hourly
where average > 80
order by
  cacheclusterid,
  timestamp;
```

### CPU hourly average < 2%

```sql
select
  cacheclusterid,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_elasticache_redis_metric_engine_cpu_utilization_hourly
where average < 2
order by
  cacheclusterid,
  timestamp;
```