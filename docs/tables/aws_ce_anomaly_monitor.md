---
title: "Steampipe Table: aws_ce_anomaly_monitor - Query AWS Cost Explorer Anomaly Monitors using SQL"
description: "Allows users to query Cost Explorer Anomaly Monitors, providing information about anomaly detection configurations, monitor status, and evaluation dates. Useful for cost governance and anomaly detection management."
folder: "Cost Explorer"
---

# Table: aws_ce_anomaly_monitor - Query AWS Cost Explorer Anomaly Monitors using SQL

AWS Cost Anomaly Detection automatically monitors your AWS spending and alerts you when spending patterns are unusual. This helps you quickly identify and investigate unexpected cost spikes using machine learning models.

## Table Usage Guide

The `aws_ce_anomaly_monitor` table in Steampipe provides you with information about Cost Explorer Anomaly Monitors in your AWS account. This table allows you to query monitor-specific details, including monitor configuration, type, dimensions, evaluation status, and metadata. You can utilize this table to gather insights on your anomaly detection setup, monitor evaluation history, and anomaly detection dimensions. The schema outlines the various attributes of an anomaly monitor including the monitor ARN, name, type, and evaluation dates.

## Examples

### List all anomaly monitors
Retrieve all anomaly monitors configured in your AWS account to review your cost anomaly detection setup.

```sql+postgres
select
  monitor_arn,
  monitor_name,
  monitor_type,
  monitor_dimension,
  creation_date,
  last_evaluated_date
from
  aws_ce_anomaly_monitor;
```

```sql+sqlite
select
  monitor_arn,
  monitor_name,
  monitor_type,
  monitor_dimension,
  creation_date,
  last_evaluated_date
from
  aws_ce_anomaly_monitor;
```

### Get monitor details with specification
View detailed information about anomaly monitors including their monitoring specifications and dimensions.

```sql+postgres
select
  monitor_arn,
  monitor_name,
  monitor_type,
  monitor_specification,
  dimensional_value_count,
  last_evaluated_date,
  last_updated_date
from
  aws_ce_anomaly_monitor;
```

```sql+sqlite
select
  monitor_arn,
  monitor_name,
  monitor_type,
  monitor_specification,
  dimensional_value_count,
  last_evaluated_date,
  last_updated_date
from
  aws_ce_anomaly_monitor;
```

### Find monitors that have not been evaluated recently
Identify anomaly monitors that haven't been evaluated in the last 7 days to ensure your monitors are actively running.

```sql+postgres
select
  monitor_name,
  monitor_type,
  creation_date,
  last_evaluated_date,
  monitor_dimension
from
  aws_ce_anomaly_monitor
where
  last_evaluated_date < current_date - interval '7 days'
order by
  last_evaluated_date;
```

```sql+sqlite
select
  monitor_name,
  monitor_type,
  creation_date,
  last_evaluated_date,
  monitor_dimension
from
  aws_ce_anomaly_monitor
where
  last_evaluated_date < date('now', '-7 days')
order by
  last_evaluated_date;
```

### Get monitor by ARN
Retrieve a specific anomaly monitor using its ARN.

```sql+postgres
select
  monitor_arn,
  monitor_name,
  monitor_type,
  creation_date,
  last_updated_date,
  dimensional_value_count
from
  aws_ce_anomaly_monitor
where
  monitor_arn = 'arn:aws:ce::123456789012:anomalymonitor/my-monitor';
```

```sql+sqlite
select
  monitor_arn,
  monitor_name,
  monitor_type,
  creation_date,
  last_updated_date,
  dimensional_value_count
from
  aws_ce_anomaly_monitor
where
  monitor_arn = 'arn:aws:ce::123456789012:anomalymonitor/my-monitor';
```

### List monitors by type
Retrieve anomaly monitors filtered by monitor type (DIMENSIONAL or CUSTOM).

```sql+postgres
select
  monitor_name,
  monitor_type,
  monitor_dimension,
  creation_date
from
  aws_ce_anomaly_monitor
where
  monitor_type = 'DIMENSIONAL'
order by
  creation_date desc;
```

```sql+sqlite
select
  monitor_name,
  monitor_type,
  monitor_dimension,
  creation_date
from
  aws_ce_anomaly_monitor
where
  monitor_type = 'DIMENSIONAL'
order by
  creation_date desc;
```

### View complete monitor specifications
Get all details including the monitor specification which contains cost categories, tags, or dimensions configuration.

```sql+postgres
select
  monitor_name,
  monitor_type,
  monitor_specification,
  dimensional_value_count,
  creation_date
from
  aws_ce_anomaly_monitor;
```

```sql+sqlite
select
  monitor_name,
  monitor_type,
  monitor_specification,
  dimensional_value_count,
  creation_date
from
  aws_ce_anomaly_monitor;
```

### Find monitors by dimension
Retrieve monitors configured to track a specific dimension for anomaly detection.

```sql+postgres
select
  monitor_name,
  monitor_type,
  monitor_dimension,
  dimensional_value_count
from
  aws_ce_anomaly_monitor
where
  monitor_dimension = 'SERVICE'
order by
  dimensional_value_count desc;
```

```sql+sqlite
select
  monitor_name,
  monitor_type,
  monitor_dimension,
  dimensional_value_count
from
  aws_ce_anomaly_monitor
where
  monitor_dimension = 'SERVICE'
order by
  dimensional_value_count desc;
```

### Get recently updated monitors
List anomaly monitors sorted by their most recent update to track recent configuration changes.

```sql+postgres
select
  monitor_name,
  monitor_type,
  last_updated_date,
  last_evaluated_date,
  creation_date
from
  aws_ce_anomaly_monitor
order by
  last_updated_date desc;
```

```sql+sqlite
select
  monitor_name,
  monitor_type,
  last_updated_date,
  last_evaluated_date,
  creation_date
from
  aws_ce_anomaly_monitor
order by
  last_updated_date desc;
```

