---
title: "Steampipe Table: aws_cost_anomaly_detection - Query AWS Cost Anomaly Detection using SQL"
description: "Allows users to query Cost Anomaly Detection monitors, providing information about anomaly detection configurations, monitor status, and evaluation dates. Useful for cost governance and optimization."
folder: "Cost Management"
---

# Table: aws_cost_anomaly_detection - Query AWS Cost Anomaly Detection using SQL

AWS Cost Anomaly Detection automatically monitors your AWS spending and alerts you when spending patterns are unusual. This helps you quickly identify and investigate unexpected cost spikes.

## Table Usage Guide

The `aws_cost_anomaly_detection` table in Steampipe provides you with information about Cost Anomaly Detection monitors in your AWS account. This table allows you to query monitor-specific details, including monitor status, frequency, detection specifications, and evaluation dates. You can utilize this table to gather insights on active anomaly monitors, their configurations, and their evaluation history. The schema outlines the various attributes of an anomaly detection monitor including the monitor ARN, name, status, and frequency.

## Examples

### List all anomaly detection monitors

```sql+postgres
select
  name,
  status,
  frequency,
  monitor_arn,
  creation_date
from
  aws_ce_cost_anomaly_detection;
```

```sql+sqlite
select
  name,
  status,
  frequency,
  monitor_arn,
  creation_date
from
  aws_ce_cost_anomaly_detection;
```

### List active anomaly detection monitors

```sql+postgres
select
  name,
  status,
  frequency,
  creation_date,
  last_evaluation_date
from
  aws_ce_cost_anomaly_detection
where
  status = 'ACTIVE';
```

```sql+sqlite
select
  name,
  status,
  frequency,
  creation_date,
  last_evaluation_date
from
  aws_ce_cost_anomaly_detection
where
  status = 'ACTIVE';
```

### Get monitor details with detection specification

```sql+postgres
select
  name,
  status,
  frequency,
  monitor_specification,
  last_evaluation_date,
  next_evaluation_date
from
  aws_ce_cost_anomaly_detection;
```

```sql+sqlite
select
  name,
  status,
  frequency,
  monitor_specification,
  last_evaluation_date,
  next_evaluation_date
from
  aws_ce_cost_anomaly_detection;
```

### Find monitors that have not been evaluated recently

```sql+postgres
select
  name,
  status,
  creation_date,
  last_evaluation_date,
  next_evaluation_date
from
  aws_ce_cost_anomaly_detection
where
  last_evaluation_date < current_date - interval '7 days';
```

```sql+sqlite
select
  name,
  status,
  creation_date,
  last_evaluation_date,
  next_evaluation_date
from
  aws_ce_cost_anomaly_detection
where
  last_evaluation_date < date('now', '-7 days');
```

### List recently created anomaly monitors

```sql+postgres
select
  name,
  status,
  frequency,
  creation_date,
  last_modified_date
from
  aws_ce_cost_anomaly_detection
where
  creation_date > current_date - interval '30 days'
order by
  creation_date desc;
```

```sql+sqlite
select
  name,
  status,
  frequency,
  creation_date,
  last_modified_date
from
  aws_ce_cost_anomaly_detection
where
  creation_date > date('now', '-30 days')
order by
  creation_date desc;
```

