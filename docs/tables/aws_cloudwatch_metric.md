# Table: aws_cloudwatch_metric

Metrics are data about the performance of your systems. By default, many services provide free metrics for resources (such as Amazon EC2 instances, Amazon EBS volumes, and Amazon RDS DB instances).

**Note**: Up to 10 dimensions can be included in the `dimensions_filter` column.

## Examples

### Basic info

```sql
select
  metric_name,
  namespace,
  dimensions
from
  aws_cloudwatch_metric;
```

### List EBS metrics

```sql
select
  metric_name,
  namespace,
  dimensions
from
  aws_cloudwatch_metric
where
  namespace = 'AWS/EBS';
```

### List EBS `VolumeReadOps` metrics

```sql
select
  metric_name,
  namespace,
  dimensions
from
  aws_cloudwatch_metric
where
  namespace = 'AWS/EBS'
  and metric_name = 'VolumeReadOps';
```

### List metrics for a specific Redshift cluster

```sql
select
  metric_name,
  namespace,
  dimensions
from
  aws_cloudwatch_metric
where
  dimensions_filter = '[
    {"Name": "ClusterIdentifier", "Value": "my-cluster-1"}
  ]'::jsonb;
```

### List EC2 API metrics

```sql
select
  metric_name,
  namespace,
  dimensions
from
  aws_cloudwatch_metric
where
  dimensions_filter = '[
    {"Name": "Type", "Value": "API"},
    {"Name": "Service", "Value": "EC2"}
  ]'::jsonb;
```
