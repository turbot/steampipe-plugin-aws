# Table: aws_iot_fleet_metric

AWS IoT Fleet Metrics is a feature in AWS IoT Device Management that allows you to monitor and manage the status and performance of your IoT device fleet. It enables you to define custom metrics based on device data reported to AWS IoT Core, such as device states or telemetry data. These metrics provide insights into various aspects of your IoT fleet, helping you make informed decisions about device maintenance, operations, and performance optimization.

## Examples

### Basic info

```sql
select
  metric_name,
  arn,
  index_name,
  creation_date,
  last_modified_date
from
  aws_iot_fleet_metric;
```

### Aggregate fleet metric by type name

```sql
select
  metric_name,
  aggregation_field,
  creation_date,
  aggregation_type_name,
  query_string
from
  aws_iot_fleet_metric
group by
  aggregation_type_name;
```

### List fleet metrics updated in the last 30 days

```sql
select
  metric_name,
  index_name,
  creation_date,
  last_modified_date,
  query_version,
  version
from
  aws_iot_fleet_metric
where
  last_modified_date >= now() - interval '30' day;
```