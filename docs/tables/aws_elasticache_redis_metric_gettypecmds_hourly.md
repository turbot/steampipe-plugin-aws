# Table: aws_elasticache_redis_metric_gettypecmds_hourly

Amazon CloudWatch Metrics provide data about the performance of your systems.  The `aws_elasticache_redis_metric_gettypecmds_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.


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
  aws_elasticache_redis_metric_gettypecmds_hourly
order by
  cacheclusterid,
  timestamp;
```



### gettypecmds sum 0ver 100 

```sql
select
  cacheclusterid,
  timestamp,
  round(minimum::numeric,2) as min_gettypecmds,
  round(maximum::numeric,2) as max_gettypecmds,
  round(average::numeric,2) as avg_gettypecmds,
  round(sum::numeric,2) as sum_gettypecmds
from
  aws_elasticache_redis_metric_gettypecmds_hourly
where sum > 100
order by
  cacheclusterid,
  timestamp;
```
