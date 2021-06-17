# Table: aws_ebs_volume_metric_write_ops_hourly

Amazon CloudWatch Metrics provide data about the performance of your systems.  The `aws_ebs_volume_metric_write_ops_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.


## Examples

### Basic info

```sql
select
  volume_id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_ebs_volume_metric_write_ops_hourly
order by
  volume_id,
  timestamp;
```

### Intervals where volumes exceed 1000 average write ops
```sql
select
  volume_id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_ebs_volume_metric_write_ops_hourly
where
  average > 1000
order by
  volume_id,
  timestamp;
```


### Intervals where volumes exceed 8000 max write ops
```sql
select
  volume_id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_ebs_volume_metric_write_ops_hourly
where
  maximum > 8000
order by
  volume_id,
  timestamp;
```



### Intervals where volume average iops exceeds provisioned iops
```sql
select 
  r.volume_id,
  r.timestamp,
  v.iops as provisioned_iops,
  round(r.average) +round(w.average) as iops_avg,
  round(r.average) as read_ops_avg,
  round(w.average) as write_ops_avg
from 
  aws_ebs_volume_metric_read_ops_hourly as r,
  aws_ebs_volume_metric_write_ops_hourly as w,
  aws_ebs_volume as v
where 
  r.volume_id = w.volume_id
  and r.timestamp = w.timestamp
  and v.volume_id = r.volume_id 
  and r.average + w.average > v.iops
order by
  r.volume_id,
  r.timestamp;
```


### Read, Write, and Total IOPS

```sql
select 
  r.volume_id,
  r.timestamp,
  round(r.average) + round(w.average) as iops_avg,
  round(r.average) as read_ops_avg,
  round(w.average) as write_ops_avg,
  round(r.maximum) + round(w.maximum) as iops_max,
  round(r.maximum) as read_ops_max,
  round(w.maximum) as write_ops_max,
  round(r.minimum) + round(w.minimum) as iops_min,
  round(r.minimum) as read_ops_min,
  round(w.minimum) as write_ops_min
from 
  aws_ebs_volume_metric_read_ops_hourly as r,
  aws_ebs_volume_metric_write_ops_hourly as w
where 
  r.volume_id = w.volume_id
  and r.timestamp = w.timestamp
order by
  r.volume_id,
  r.timestamp;
```