---
title: "Steampipe Table: aws_application_signals_service_level_objective - Query AWS CloudWatch Application Signals Service Level Objectives using SQL"
description: "Allows users to query AWS CloudWatch Application Signals for information about their service level objectives."
folder: "Application Signals"
---

# Table: aws_application_signals_service_level_objective - Query AWS CloudWatch Application Signals Service Level Objectives using SQL

The AWS CloudWatch Application Signals is a service that allows you to monitor and improve application performance.

## Table Usage Guide

The `aws_application_signals_service_level_objective` table in Steampipe provides you with information about CloudWatch
Application Signals service level objectives. This table allows you, as a DevOps engineer, to configure service level
objectives based on a variety of metrics and monitor your applications' performance. You can utilize this table to
manage and monitor service level objectives of your applications. The schema outlines the various attributes of the
service level objectives for you, including their metric sources and alarm thresholds.

## Examples

### Fetch Service Level Objectives Based on CloudWatch Metrics
Fetch service level objectives configured based on CloudWatch metrics.

```sql+postgres
select
  name,
  burn_rate_configurations,
  goal
from
  aws_application_signals_slo
where
  metric_source_type = 'CloudWatchMetric'
```

```sql+sqlite
select
  name,
  burn_rate_configurations,
  goal
from
  aws_application_signals_slo
where
  metric_source_type = 'CloudWatchMetric'
```
