# Table: aws_rds_db_instance_metric_connections_hourly

Amazon CloudWatch Metrics provide data about the performance of your systems.  The `aws_rds_db_instance_metric_connections_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.


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
  aws_rds_db_instance_metric_connections_hourly
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
  aws_rds_db_instance_metric_connections_hourly
where 
  average > 100
order by
  db_instance_identifier,
  timestamp;
```