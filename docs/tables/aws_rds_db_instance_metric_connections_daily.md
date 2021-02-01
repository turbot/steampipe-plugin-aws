# Table: aws_rds_db_instance_metric_connections_daily

Amazon CloudWatch Metrics provide data about the performance of your systems.  The `aws_rds_db_instance_metric_connections_daily` table provides metric statistics at 24 hour intervals for the past year.


## Examples

### Basic info

```sql
select
  db_instance_identifier,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_rds_db_instance_metric_connections_daily
order by
  db_instance_identifier,
  timestamp;
```


### Intervals averaging over 100 connections

```sql
select
  db_instance_identifier,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_rds_db_instance_metric_connections_daily
where 
  average > 100
order by
  db_instance_identifier,
  timestamp;
```


### Instances with no connections in the past week

```sql
select
  db_instance_identifier,
  sum(maximum) as total_connections
from
  aws_rds_db_instance_metric_connections
where 
  timestamp > (current_date - interval '7' day)
group by
  db_instance_identifier
having
  sum(maximum) = 0 
;
```


