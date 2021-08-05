# Table: aws_emr_cluster_metric_isidle

Amazon CloudWatch Metrics provide data about the performance of your systems.  The `aws_emr_cluster_metric_is_idle` table provides metric statistics at 5 minute intervals for the most recent 5 days.

## Examples

### Basic info

```sql
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_emr_cluster_metric_is_idle
order by
  id,
  timestamp;
```