# Table: aws_elasticache_redis_metric_new_connections_hourly

Amazon CloudWatch Metrics provide data about the performance of your systems.The `aws_elasticache_redis_metric_new_connections_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

## Examples

### Basic info

```sql
select
  cacheclusterid,
  timestamp,
  minimum,
  maximum,
  average
from
  aws_elasticache_redis_metric_new_connections_hourly
order by
  cacheclusterid,
  timestamp;
```

### newconnections sum over 10

```sql
select
  cacheclusterid,
  timestamp,
  round(minimum::numeric,2) as min_newconnections,
  round(maximum::numeric,2) as max_newconnections,
  round(average::numeric,2) as avg_newconnections,
  round(sum::numeric,2) as sum_newconnections
from
  aws_elasticache_redis_metric_new_connections_hourly
where sum > 10
order by
  cacheclusterid,
  timestamp;
```