---
title: "Steampipe Table: aws_application_signals_service_level_objective - Query AWS Application Signals using SQL"
description: "Allows users to query AWS Application Signals to retrieve detailed information about each SLO, including its name, ARN, attainment_goal, goal and sli."
---

# Table: aws_application_signals_service_level_objective - Query AWS Application Signals Service Level Objective using SQL

## Table Usage Guide

The `aws_application_signals_service_level_objective` table in Steampipe provides you with information about SLOs within AWS Application Signals. This table allows you to query SLO details, including the name, ARN, attainment_goal, goal and sli, and associated metadata.

## Examples

```sql+postgres
select
  arn,
  name,
  attainment_goal,
  goal,
from
  aws_application_signals_service_level_objective
```

```sql+postgres
select
  arn,
  name,
  sli::json -> 'ComparisonOperator' as "Must Be",
  sli::json -> 'MetricThreshold' as "Threshold"
  region
from
  aws_application_signals_service_level_objective
```
