# Table: aws_elasticache_redis_metric_list_based_cmds_hourly

Amazon CloudWatch Metrics provide data about the performance of your systems. The `aws_elasticache_redis_metric_list_based_cmds_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

## Examples

### Basic info

```sql
select
  cache_cluster_id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count,
  sum
from
  aws_elasticache_redis_metric_list_based_cmds_hourly
order by
  cache_cluster_id,
  timestamp;
```

### listbasedcmds sum over 100 

```sql
select
  cache_cluster_id,
  timestamp,
  round(minimum::numeric,2) as min_listbasedcmds,
  round(maximum::numeric,2) as max_listbasedcmds,
  round(average::numeric,2) as avg_listbasedcmds,
  round(sum::numeric,2) as sum_listbasedcmds
from
  aws_elasticache_redis_metric_list_based_cmds_hourly
where sum > 100
order by
  cache_cluster_id,
  timestamp;
```