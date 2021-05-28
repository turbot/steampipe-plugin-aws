# Table: aws_rds_db_instance_metric_cpu_utilization_daily

Amazon CloudWatch Metrics provide data about the performance of your systems.  The `aws_rds_db_instance_metric_cpu_utilization_daily` table provides metric statistics at 24 hour intervals for the last year.


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
  aws_rds_db_instance_metric_cpu_utilization_daily
order by
  db_instance_identifier,
  timestamp;
```



### CPU Over 80% average

```sql
select
  db_instance_identifier,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_rds_db_instance_metric_cpu_utilization_daily
where average > 80
order by
  db_instance_identifier,
  timestamp;
```

### CPU daily average < 2%

```sql
select
  db_instance_identifier,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_rds_db_instance_metric_cpu_utilization_daily
where average < 2
order by
  db_instance_identifier,
  timestamp;
```