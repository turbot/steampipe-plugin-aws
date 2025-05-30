---
title: "Steampipe Table: aws_cost_by_resource_daily - Query AWS Cost Explorer Resource Costs using SQL"
description: "Allows users to query AWS Cost Explorer Resource Costs on a daily basis, providing detailed cost information for individual AWS resources."
folder: "Cost Explorer"
---

# Table: aws_cost_by_resource_daily - Query AWS Cost Explorer Resource Costs using SQL

AWS Cost Explorer helps you visualize, understand, and manage your AWS costs and usage. The Resource Cost feature provides detailed cost information at the individual resource level, helping you identify specific resources that drive your AWS costs.

## Table Usage Guide

The `aws_cost_by_resource_daily` table provides insights into resource-level costs within AWS Cost Explorer. This table allows you, as a financial analyst or cloud administrator, to query daily cost details for specific AWS resources, helping you understand spending patterns and identify cost optimization opportunities. The schema outlines various cost metrics including unblended cost, amortized cost, and usage quantity, along with resource identifiers and time periods.

**Important Notes**

- The [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request you make will incur a cost of $0.01.
- By default, the table shows resource-level data with `dimension_key = 'LINKED_ACCOUNT'` and `dimension_value` set to the caller's AWS account ID for the last 14 days. For historical data beyond 14 days, you need to enable the cost allocation data in your AWS Cost Explorer settings.
- This table supports optional quals. Queries with optional quals are optimised to reduce query time and cost. Optional quals are supported for the following columns:
  - `resource_id` with supported operators `=` and `<>`.
  - `dimension_key` with supported operator `=`.
  - `dimension_value` with supported operator `=`.
  - `period_start` with supported operators `=`, `>=`, `>`, `<=`, and `<`.
  - `period_end` with supported operators `=`, `>=`, `>`, `<=`, and `<`.

## Examples

### Basic info
Get a simple overview of daily resource costs with essential fields.

```sql+postgres
select
  resource_id,
  period_start,
  period_end,
  dimension_key,
  dimension_value,
  blended_cost_amount::numeric::money as blended_cost,
  unblended_cost_amount::numeric::money as unblended_cost
from
  aws_cost_by_resource_daily
where
  period_start >= current_date - interval '14 days'
order by
  period_start desc;
```

```sql+sqlite
select
  resource_id,
  period_start,
  period_end,
  dimension_key,
  dimension_value,
  cast(blended_cost_amount as decimal) as blended_cost,
  cast(unblended_cost_amount as decimal) as unblended_cost
from
  aws_cost_by_resource_daily
where
  period_start >= date('now', '-14 days')
order by
  period_start desc;
```

### Daily cost for a specific EC2 instance in a region

Explore the daily costs associated with a particular EC2 instance in a specific region to track its financial impact over time.

```sql+postgres
select
  resource_id,
  period_start,
  period_end,
  dimension_key,
  dimension_value,
  blended_cost_amount::numeric::money as blended_cost,
  unblended_cost_amount::numeric::money as unblended_cost,
  net_unblended_cost_amount::numeric::money as net_unblended_cost,
  amortized_cost_amount::numeric::money as amortized_cost,
  net_amortized_cost_amount::numeric::money as net_amortized_cost,
  usage_quantity_amount,
  usage_quantity_unit,
  normalized_usage_amount,
  normalized_usage_unit
from
  aws_cost_by_resource_daily
where
  resource_id = 'i-1234567890abcdef0'
  and dimension_key = 'REGION'
  and dimension_value = 'us-east-1'
  and period_start >= current_date - interval '14 days'
order by
  period_start desc;
```

```sql+sqlite
select
  resource_id,
  period_start,
  period_end,
  dimension_key,
  dimension_value,
  cast(blended_cost_amount as decimal) as blended_cost,
  cast(unblended_cost_amount as decimal) as unblended_cost,
  cast(net_unblended_cost_amount as decimal) as net_unblended_cost,
  cast(amortized_cost_amount as decimal) as amortized_cost,
  cast(net_amortized_cost_amount as decimal) as net_amortized_cost,
  usage_quantity_amount,
  usage_quantity_unit,
  normalized_usage_amount,
  normalized_usage_unit
from
  aws_cost_by_resource_daily
where
  resource_id = 'i-1234567890abcdef0'
  and dimension_key = 'REGION'
  and dimension_value = 'us-east-1'
  and period_start >= date('now', '-14 days')
order by
  period_start desc;
```

### Top 10 most expensive resources yesterday by region

Identify the resources that generated the highest costs yesterday in a specific region to focus cost optimization efforts.

```sql+postgres
select
  resource_id,
  dimension_key,
  dimension_value,
  blended_cost_amount::numeric::money as blended_cost,
  unblended_cost_amount::numeric::money as unblended_cost,
  net_unblended_cost_amount::numeric::money as net_unblended_cost,
  amortized_cost_amount::numeric::money as amortized_cost,
  net_amortized_cost_amount::numeric::money as net_amortized_cost,
  usage_quantity_amount,
  usage_quantity_unit,
  normalized_usage_amount,
  normalized_usage_unit
from
  aws_cost_by_resource_daily
where
  period_start = date_trunc('day', current_timestamp - interval '1 day')
  and dimension_key = 'REGION'
  and dimension_value = 'us-east-1'
order by
  blended_cost_amount desc
limit 10;
```

```sql+sqlite
select
  resource_id,
  dimension_key,
  dimension_value,
  cast(blended_cost_amount as decimal) as blended_cost,
  cast(unblended_cost_amount as decimal) as unblended_cost,
  cast(net_unblended_cost_amount as decimal) as net_unblended_cost,
  cast(amortized_cost_amount as decimal) as amortized_cost,
  cast(net_amortized_cost_amount as decimal) as net_amortized_cost,
  usage_quantity_amount,
  usage_quantity_unit,
  normalized_usage_amount,
  normalized_usage_unit
from
  aws_cost_by_resource_daily
where
  period_start = date('now', '-1 day')
  and dimension_key = 'REGION'
  and dimension_value = 'us-east-1'
order by
  blended_cost_amount desc
limit 10;
```

### Daily costs for a specific resource by linked account

Analyze how costs for a particular resource are distributed across different linked accounts.

```sql+postgres
select
  resource_id,
  period_start,
  dimension_key,
  dimension_value,
  blended_cost_amount::numeric::money as blended_cost,
  net_unblended_cost_amount::numeric::money as net_unblended_cost,
  usage_quantity_amount,
  usage_quantity_unit,
  normalized_usage_amount,
  normalized_usage_unit
from
  aws_cost_by_resource_daily
where
  resource_id = 'arn:aws:rds:us-east-1:123456789012:db:my-database'
  and dimension_key = 'LINKED_ACCOUNT'
  and dimension_value = '123456789012'
  and period_start >= date_trunc('month', current_date)
order by
  period_start desc,
  blended_cost_amount desc;
```

```sql+sqlite
select
  resource_id,
  period_start,
  dimension_key,
  dimension_value,
  cast(blended_cost_amount as decimal) as blended_cost,
  cast(net_unblended_cost_amount as decimal) as net_unblended_cost,
  usage_quantity_amount,
  usage_quantity_unit,
  normalized_usage_amount,
  normalized_usage_unit
from
  aws_cost_by_resource_daily
where
  resource_id = 'arn:aws:rds:us-east-1:123456789012:db:my-database'
  and dimension_key = 'LINKED_ACCOUNT'
  and dimension_value = '123456789012'
  and period_start >= date('now', 'start of month')
order by
  period_start desc,
  blended_cost_amount desc;
```

### Resources with unusual cost patterns by service
Identify resources that have experienced significant cost increases compared to their previous day within a specific service.

```sql+postgres
with daily_costs as (
  select
    resource_id,
    dimension_key,
    dimension_value,
    period_start,
    blended_cost_amount::numeric as cost,
    net_unblended_cost_amount::numeric as net_cost,
    usage_quantity_amount,
    usage_quantity_unit,
    normalized_usage_amount,
    normalized_usage_unit,
    lag(blended_cost_amount::numeric) over (partition by resource_id order by period_start) as previous_day_cost
  from
    aws_cost_by_resource_daily
  where
    period_start >= date_trunc('month', current_date)
    and dimension_key = 'SERVICE'
    and dimension_value = 'Amazon Elastic Compute Cloud - Compute'
)
select
  resource_id,
  dimension_key,
  dimension_value,
  period_start,
  cost::money as cost,
  net_cost::money as net_cost,
  usage_quantity_amount,
  usage_quantity_unit,
  normalized_usage_amount,
  normalized_usage_unit,
  previous_day_cost::money as previous_day_cost,
  ((cost - previous_day_cost) / nullif(previous_day_cost, 0) * 100)::numeric(10,2) as cost_increase_percent
from
  daily_costs
where
  previous_day_cost > 0
  and cost > previous_day_cost * 2
order by
  cost_increase_percent desc;
```

```sql+sqlite
with daily_costs as (
  select
    resource_id,
    dimension_key,
    dimension_value,
    period_start,
    cast(blended_cost_amount as decimal) as cost,
    cast(net_unblended_cost_amount as decimal) as net_cost,
    usage_quantity_amount,
    usage_quantity_unit,
    normalized_usage_amount,
    normalized_usage_unit,
    lag(cast(blended_cost_amount as decimal)) over (partition by resource_id order by period_start) as previous_day_cost
  from
    aws_cost_by_resource_daily
  where
    period_start >= date('now', 'start of month')
    and dimension_key = 'SERVICE'
    and dimension_value = 'Amazon Elastic Compute Cloud - Compute'
)
select
  resource_id,
  dimension_key,
  dimension_value,
  period_start,
  cost,
  net_cost,
  usage_quantity_amount,
  usage_quantity_unit,
  normalized_usage_amount,
  normalized_usage_unit,
  previous_day_cost,
  round(((cost - previous_day_cost) / case when previous_day_cost = 0 then null else previous_day_cost end * 100), 2) as cost_increase_percent
from
  daily_costs
where
  previous_day_cost > 0
  and cost > previous_day_cost * 2
order by
  cost_increase_percent desc;
```
