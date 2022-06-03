# Table: aws_cloudwatch_metric

Metrics are data about the performance of your systems. By default, many services provide free metrics for resources (such as Amazon EC2 instances, Amazon EBS volumes, and Amazon RDS DB instances).

## Examples

### Basic info

```sql
select
  name,
  namespace,
  dimension_name,
  dimension_value
from
  aws_cloudwatch_metric;
```

### List metric by EBS namespace

```sql
select
  name,
  namespace,
  dimension_name,
  dimension_value
from
  aws_cloudwatch_metric
where
  namespace = 'AWS/EBS';
```

### List metric details for metric name VolumeReadOps

```sql
select
  name,
  namespace,
  dimension_name,
  dimension_value
from
  aws_cloudwatch_metric
where
  name = 'VolumeReadOps';
```

### List metric for a redshift cluster

```sql
select
  name,
  namespace,
  dimension_name,
  dimension_value
from
  aws_cloudwatch_metric
where
  dimension_name = 'ClusterIdentifier' and dimension_value = 'redshift-cluster-1';
```