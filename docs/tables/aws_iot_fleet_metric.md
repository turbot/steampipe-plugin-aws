---
title: "Steampipe Table: aws_iot_fleet_metric - Query AWS IoT Fleet Metrics using SQL"
description: "Allows users to query AWS IoT Fleet Metrics to gain insights into each fleet metric's configuration, including ARN, creation date, and aggregation information."
folder: "IoT Core"
---

# Table: aws_iot_fleet_metric - Query AWS IoT Fleet Metrics using SQL

AWS IoT Fleet Metrics, part of AWS IoT Device Management, allows for the monitoring and management of the status and performance of your IoT device fleet. It enables the definition of custom metrics based on device data reported to AWS IoT Core, such as device states or telemetry data. These metrics are instrumental in providing insights into various aspects of your IoT fleet, aiding in decisions related to device maintenance, operations, and performance optimization.

## Table Usage Guide

The `aws_iot_fleet_metric` table can be utilized to access detailed information about custom metrics defined for IoT device fleets. This table is essential for IoT administrators and analysts who need to oversee fleet performance, status, and operational metrics within AWS.

## Examples

### Basic info
Retrieve fundamental details about AWS IoT Fleet Metrics, including metric names, associated ARNs, index names, and dates of creation and modification.

```sql+postgres
select
  metric_name,
  arn,
  index_name,
  creation_date,
  last_modified_date
from
  aws_iot_fleet_metric;
```

```sql+sqlite
select
  metric_name,
  arn,
  index_name,
  creation_date,
  last_modified_date
from
  aws_iot_fleet_metric;
```

### Group fleet metrics by aggregation type name
Group fleet metrics by their aggregation type name. This query is useful for analyzing metrics across different aggregation types, providing a broader view of fleet data categorization.

```sql+postgres
select
  metric_name,
  aggregation_type_name
from
  aws_iot_fleet_metric
group by
  aggregation_type_name,
  metric_name;
```

```sql+sqlite
select
  metric_name,
  aggregation_type_name
from
  aws_iot_fleet_metric
group by
  aggregation_type_name,
  metric_name;
```

### List fleet metrics updated in the last 30 days
Find fleet metrics that have been updated within the last 30 days. This query assists in identifying recent changes or updates in fleet metric configurations.

```sql+postgres
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
  last_modified_date >= now() - interval '30 days';
```

```sql+sqlite
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
  datetime(last_modified_date) >= datetime('now', '-30 days');
```
