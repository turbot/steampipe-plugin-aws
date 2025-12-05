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
Retrieve all budgets in your AWS account to get an overview of your cost management setup.

```sql+postgres
select
  name,
  type,
  limit_amount,
  time_unit
from
  aws_budgets_budget;
```

```sql+sqlite
select
  name,
  type,
  limit_amount,
  time_unit
from
  aws_budgets_budget;
```

### Get budget spending details
View detailed spending information for each budget, including actual spending and forecasted amounts.

```sql+postgres
select
  name,
  limit_amount,
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
  limit_amount,
  calculated_spend_actual_spend,
  calculated_spend_forecasted_spend,
  time_period_start,
  time_period_end
from
  aws_budgets_budget;
```

### List budgets with notifications
Find all budgets that have associated notifications configured for cost alerts.

```sql+postgres
select
  name,
  type,
  limit_amount,
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
  limit_amount,
  notifications
from
  aws_budgets_budget
where
  notifications is not null;
```

### Find budgets approaching their limit
Identify budgets that have reached 80% or more of their configured limit to enable proactive cost management.

```sql+postgres
select
  name,
  type,
  limit_amount,
  calculated_spend_actual_spend,
  (cast(calculated_spend_actual_spend as float) / cast(limit_amount as float) * 100) as percent_of_budget
from
  aws_budgets_budget
where
  cast(calculated_spend_actual_spend as float) > (cast(limit_amount as float) * 0.8);
```

```sql+sqlite
select
  name,
  type,
  limit_amount,
  calculated_spend_actual_spend,
  (cast(calculated_spend_actual_spend as float) / cast(limit_amount as float) * 100) as percent_of_budget
from
  aws_budgets_budget
where
  cast(calculated_spend_actual_spend as float) > (cast(limit_amount as float) * 0.8);
```

### List monthly budgets
Query all budgets configured with a monthly time unit for tracking recurring costs.

```sql+postgres
select
  name,
  type,
  limit_amount,
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
  limit_amount,
  time_unit,
  time_period_start,
  time_period_end
from
  aws_budgets_budget
where
  time_unit = 'MONTHLY';
```

### View cost filters
Display all budgets with their cost filter settings to understand how budgets are scoped to specific services or dimensions.

```sql+postgres
select
  name,
  type,
  limit_amount,
  cost_filters
from
  aws_budgets_budget;
```

```sql+sqlite
select
  name,
  type,
  limit_amount,
  cost_filters
from
  aws_budgets_budget;
```

### View cost types
Show all budgets and the cost types they include (such as taxes, support, recurring charges, upfront costs, etc.).

```sql+postgres
select
  name,
  type,
  limit_amount,
  cost_types
from
  aws_budgets_budget;
```

```sql+sqlite
select
  name,
  type,
  limit_amount,
  cost_types
from
  aws_budgets_budget;
```

### Filter budgets by service using cost_filters
Query budgets that are filtered to monitor specific AWS services.

```sql+postgres
select
  name,
  type,
  limit_amount,
  cost_filters
from
  aws_budgets_budget
where
  cost_filters is not null;
```

```sql+sqlite
select
  name,
  type,
  limit_amount,
  cost_filters
from
  aws_budgets_budget
where
  cost_filters is not null;
```

### View cost type breakdown
Examine detailed cost type configurations to understand what is included in budget calculations (blended vs. unblended costs, taxes, support, etc.).

```sql+postgres
select
  name,
  cost_types->'IncludeTax' as include_tax,
  cost_types->'IncludeSupport' as include_support,
  cost_types->'IncludeUpfront' as include_upfront,
  cost_types->'IncludeRecurring' as include_recurring,
  cost_types->'UseBlended' as use_blended,
  cost_types->'UseAmortized' as use_amortized
from
  aws_budgets_budget;
```

```sql+sqlite
select
  name,
  json_extract(cost_types, '$.IncludeTax') as include_tax,
  json_extract(cost_types, '$.IncludeSupport') as include_support,
  json_extract(cost_types, '$.IncludeUpfront') as include_upfront,
  json_extract(cost_types, '$.IncludeRecurring') as include_recurring,
  json_extract(cost_types, '$.UseBlended') as use_blended,
  json_extract(cost_types, '$.UseAmortized') as use_amortized
from
  aws_budgets_budget;
```

### Budgets using blended costs vs amortized costs
Compare budgets configured with different cost calculation methods to understand cost accounting approaches.

```sql+postgres
select
  name,
  limit_amount,
  cost_types->>'UseBlended' as use_blended,
  cost_types->>'UseAmortized' as use_amortized
from
  aws_budgets_budget;
```

```sql+sqlite
select
  name,
  limit_amount,
  json_extract(cost_types, '$.UseBlended') as use_blended,
  json_extract(cost_types, '$.UseAmortized') as use_amortized
from
  aws_budgets_budget;
```

### Budgets that include support costs
Find all budgets that are configured to include AWS support costs in budget calculations.

```sql+postgres
select
  name,
  limit_amount,
  cost_types->>'IncludeSupport' as include_support,
  cost_types->>'IncludeTax' as include_tax
from
  aws_budgets_budget
where
  (cost_types->>'IncludeSupport')::boolean = true;
```

```sql+sqlite
select
  name,
  limit_amount,
  json_extract(cost_types, '$.IncludeSupport') as include_support,
  json_extract(cost_types, '$.IncludeTax') as include_tax
from
  aws_budgets_budget
where
  json_extract(cost_types, '$.IncludeSupport') = 'true';
```

### Budgets with service-specific cost filters
Query budgets that filter spending by specific AWS services or service groups.

```sql+postgres
select
  name,
  type,
  limit_amount,
  cost_filters->'Service' as filtered_services
from
  aws_budgets_budget
where
  cost_filters ? 'Service';
```

```sql+sqlite
select
  name,
  type,
  limit_amount,
  json_extract(cost_filters, '$.Service') as filtered_services
from
  aws_budgets_budget
where
  json_extract(cost_filters, '$.Service') is not null;
```

### Budgets filtered by linked account
Find budgets that are scoped to specific AWS linked accounts in an organization.

```sql+postgres
select
  name,
  type,
  limit_amount,
  cost_filters->'LinkedAccount' as linked_accounts
from
  aws_budgets_budget
where
  cost_filters ? 'LinkedAccount';
```

```sql+sqlite
select
  name,
  type,
  limit_amount,
  json_extract(cost_filters, '$.LinkedAccount') as linked_accounts
from
  aws_budgets_budget
where
  json_extract(cost_filters, '$.LinkedAccount') is not null;
```

### Budgets filtered by region
Query budgets that are scoped to monitor costs in specific AWS regions.

```sql+postgres
select
  name,
  type,
  limit_amount,
  cost_filters->'Region' as regions
from
  aws_budgets_budget
where
  cost_filters ? 'Region';
```

```sql+sqlite
select
  name,
  type,
  limit_amount,
  json_extract(cost_filters, '$.Region') as regions
from
  aws_budgets_budget
where
  json_extract(cost_filters, '$.Region') is not null;
```

### Budgets that exclude refunds and credits
Find budgets that are configured without refunds or credits in their cost calculations.

```sql+postgres
select
  name,
  limit_amount,
  cost_types->>'IncludeRefund' as include_refund,
  cost_types->>'IncludeCredit' as include_credit
from
  aws_budgets_budget
where
  (cost_types->>'IncludeRefund')::boolean = false
  or (cost_types->>'IncludeCredit')::boolean = false;
```

```sql+sqlite
select
  name,
  limit_amount,
  json_extract(cost_types, '$.IncludeRefund') as include_refund,
  json_extract(cost_types, '$.IncludeCredit') as include_credit
from
  aws_budgets_budget
where
  json_extract(cost_types, '$.IncludeRefund') = 'false'
  or json_extract(cost_types, '$.IncludeCredit') = 'false';
```

### Cost types summary by budget
Get a summary view of how different cost types are configured across all budgets.

```sql+postgres
select
  name,
  (cost_types->>'IncludeTax')::boolean as include_tax,
  (cost_types->>'IncludeSupport')::boolean as include_support,
  (cost_types->>'IncludeSubscription')::boolean as include_subscription,
  (cost_types->>'IncludeRefund')::boolean as include_refund,
  (cost_types->>'UseBlended')::boolean as use_blended
from
  aws_budgets_budget
order by
  name;
```

```sql+sqlite
select
  name,
  json_extract(cost_types, '$.IncludeTax') as include_tax,
  json_extract(cost_types, '$.IncludeSupport') as include_support,
  json_extract(cost_types, '$.IncludeSubscription') as include_subscription,
  json_extract(cost_types, '$.IncludeRefund') as include_refund,
  json_extract(cost_types, '$.UseBlended') as use_blended
from
  aws_budgets_budget
order by
  name;
```


