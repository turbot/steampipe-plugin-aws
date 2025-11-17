---
title: "Steampipe Table: aws_budgets_budget - Query AWS Budgets using SQL"
description: "Allows users to query AWS Budgets data, providing detailed information about cost and usage budgets, budget limits, and spending notifications. Useful for cost governance, alerting, and optimization workflows."
folder: "Budgets"
---

# Table: aws_budgets_budget - Query AWS Budgets using SQL

AWS Budgets enables you to set custom cost and usage budgets and receive notifications when your usage exceeds the budgeted amounts. This helps you track and control spending across your AWS environment.

## Table Usage Guide

The `aws_budgets_budget` table in Steampipe provides you with information about each budget in your AWS account. This table allows you to query budget-specific details, including budget limits, spending amounts, notification settings, and associated metadata. You can utilize this table to gather insights on budget types, time periods, cost filters, and spending forecasts. The schema outlines the various attributes of a budget including the budget name, type, limit, calculated spend, and associated notifications.

## Examples

### List all budgets

```sql+postgres
select
  name,
  type,
  limit,
  time_unit,
  status
from
  aws_budgets_budget;
```

```sql+sqlite
select
  name,
  type,
  limit,
  time_unit,
  status
from
  aws_budgets_budget;
```

### Get budget spending details

```sql+postgres
select
  name,
  limit,
  calculated_spend_actual_spend,
  calculated_spend_forecasted_spend,
  time_period_start,
  time_period_end
from
  aws_budgets_budget;
```

```sql+sqlite
select
  name,
  limit,
  calculated_spend_actual_spend,
  calculated_spend_forecasted_spend,
  time_period_start,
  time_period_end
from
  aws_budgets_budget;
```

### List budgets with notifications

```sql+postgres
select
  name,
  type,
  limit,
  notifications
from
  aws_budgets_budget
where
  notifications is not null;
```

```sql+sqlite
select
  name,
  type,
  limit,
  notifications
from
  aws_budgets_budget
where
  notifications is not null;
```

### Find budgets approaching their limit

```sql+postgres
select
  name,
  type,
  limit,
  calculated_spend_actual_spend,
  (cast(calculated_spend_actual_spend as float) / cast(limit as float) * 100) as percent_of_budget
from
  aws_budgets_budget
where
  cast(calculated_spend_actual_spend as float) > (cast(limit as float) * 0.8);
```

```sql+sqlite
select
  name,
  type,
  limit,
  calculated_spend_actual_spend,
  (cast(calculated_spend_actual_spend as float) / cast(limit as float) * 100) as percent_of_budget
from
  aws_budgets_budget
where
  cast(calculated_spend_actual_spend as float) > (cast(limit as float) * 0.8);
```

### List monthly budgets

```sql+postgres
select
  name,
  type,
  limit,
  time_unit,
  time_period_start,
  time_period_end
from
  aws_budgets_budget
where
  time_unit = 'MONTHLY';
```

```sql+sqlite
select
  name,
  type,
  limit,
  time_unit,
  time_period_start,
  time_period_end
from
  aws_budgets_budget
where
  time_unit = 'MONTHLY';
```

