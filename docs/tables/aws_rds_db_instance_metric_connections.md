# Table: aws_rds_db_instance_metric_connections

Amazon CloudWatch Metrics provide data about the performance of your systems.  The `aws_rds_db_instance_metric_connections` table provides metric statistics at 5 minute intervals for the most recent 5 days.


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
  aws_rds_db_instance_metric_connections
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
  aws_rds_db_instance_metric_connections
where 
  average > 100
order by
  db_instance_identifier,
  timestamp;
```