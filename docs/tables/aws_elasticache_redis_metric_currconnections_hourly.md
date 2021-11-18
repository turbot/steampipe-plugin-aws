# Table: aws_elasticache_redis_metric_currconnections_hourly

Amazon CloudWatch Metrics provide data about the performance of your systems.The `aws_elasticache_redis_metric_currconnections_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

## Examples

### Basic info

```sql
select
  cacheclusterid,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_elasticache_redis_metric_currconnections_hourly
order by
  cacheclusterid,
  timestamp;
```

### currconnections Over 100 average

```sql
select
  cacheclusterid,
  timestamp,
  round(minimum::numeric,2) as min_currconnections,
  round(maximum::numeric,2) as max_currconnections,
  round(average::numeric,2) as avg_currconnections,
  sample_count
from
  aws_elasticache_redis_metric_currconnections_hourly
where average > 100
order by
  cacheclusterid,
  timestamp;
```

