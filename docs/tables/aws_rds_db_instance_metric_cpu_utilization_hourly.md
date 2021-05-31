# Table: aws_rds_db_instance_metric_cpu_utilization_hourly

Amazon CloudWatch Metrics provide data about the performance of your systems.  The `aws_rds_db_instance_metric_cpu_utilization_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.


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
  aws_rds_db_instance_metric_cpu_utilization_hourly
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
  aws_rds_db_instance_metric_cpu_utilization_hourly
where average > 80
order by
  db_instance_identifier,
  timestamp;
```

### CPU hourly average < 2%

```sql
select
  db_instance_identifier,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  aws_rds_db_instance_metric_cpu_utilization_hourly
where average < 2
order by
  db_instance_identifier,
  timestamp;
```