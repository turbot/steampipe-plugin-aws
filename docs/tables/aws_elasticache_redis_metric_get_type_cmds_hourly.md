# Table: aws_elasticache_redis_metric_get_type_cmds_hourly

Amazon CloudWatch Metrics provide data about the performance of your systems. The `aws_elasticache_redis_metric_get_type_cmds_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

## Examples

### Basic info

```sql
select
  cache_cluster_id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_elasticache_redis_metric_get_type_cmds_hourly
order by
  cache_cluster_id,
  timestamp;
```

### gettypecmds sum 0ver 100 

```sql
select
  cache_cluster_id,
  timestamp,
  round(minimum::numeric,2) as min_gettypecmds,
  round(maximum::numeric,2) as max_gettypecmds,
  round(average::numeric,2) as avg_gettypecmds,
  round(sum::numeric,2) as sum_gettypecmds
from
  aws_elasticache_redis_metric_get_type_cmds_hourly
where sum > 100
order by
  cache_cluster_id,
  timestamp;
```
