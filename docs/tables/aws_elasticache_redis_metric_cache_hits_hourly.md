# Table: aws_elasticache_redis_metric_cache_hits_hourly

Amazon CloudWatch Metrics provide data about the performance of your systems.The `aws_elasticache_redis_metric_cache_hits_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

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
  aws_elasticache_redis_metric_cache_hits_hourly
order by
  cacheclusterid,
  timestamp;
```


### CacheHit sum below 10 

```sql
select
  cacheclusterid,
  timestamp,
  round(sum::numeric,2) as sum_cachehits,
  round(average::numeric,2) as average_cachehits,
  sample_count
from
  aws_elasticache_redis_metric_cache_hits_hourly
where sum < 10
order by
  cacheclusterid,
  timestamp;
```

### CacheHit hourly average < 100

```sql
select
  cacheclusterid,
  timestamp,
  round(minimum::numeric,2) as min_cachehits,
  round(maximum::numeric,2) as max_cachehits,
  round(average::numeric,2) as avg_cachehits,
  sample_count
from
  aws_elasticache_redis_metric_cache_hits_hourly
where average < 100
order by
  cacheclusterid,
  timestamp;
```