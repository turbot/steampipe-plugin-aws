---
title: "Steampipe Table: aws_cost_by_resource_monthly - Query AWS Cost Explorer Resource Costs using SQL"
description: "Allows users to query AWS Cost Explorer Resource Costs on a monthly basis, providing detailed cost information for individual AWS resources."
folder: "Cost Explorer"
---

# Table: aws_cost_by_resource_monthly - Query AWS Cost Explorer Resource Costs using SQL

AWS Cost Explorer helps you visualize, understand, and manage your AWS costs and usage. The Resource Cost feature provides detailed cost information at the individual resource level with monthly granularity, helping you identify specific resources that drive your AWS costs and understand long-term cost patterns.

## Table Usage Guide

The `aws_cost_by_resource_monthly` table provides insights into resource-level costs within AWS Cost Explorer with monthly granularity. This table allows you, as a financial analyst or cloud administrator, to query monthly cost details for specific AWS resources, helping you understand long-term spending patterns and identify cost optimization opportunities. The schema outlines various cost metrics including unblended cost, amortized cost, and usage quantity, along with resource identifiers and time periods.

**Important Notes**

- The [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request you make will incur a cost of $0.01.
- By default, the table shows resource-level data with `dimension_key = 'LINKED_ACCOUNT'` and `dimension_value` set to the caller's AWS account ID for the last 14 days. For historical data beyond 14 days, you need to enable historical cost allocation data in your AWS Cost Explorer settings.
- This table supports optional quals. Queries with optional quals are optimised to reduce query time and cost. Optional quals are supported for the following columns:
  - `resource_id` with supported operators `=` and `<>`.
  - `dimension_key` with supported operator `=`.
  - `dimension_value` with supported operator `=`.
  - `period_start` with supported operators `=`, `>=`, `>`, `<=`, and `<`.
  - `period_end` with supported operators `=`, `>=`, `>`, `<=`, and `<`.

## Examples

### Basic info
Get a simple overview of resource costs with essential fields.

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
  aws_cost_by_resource_monthly
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
  aws_cost_by_resource_monthly
where
  period_start >= date('now', '-14 days')
order by
  period_start desc;
```

### Monthly cost trend for a specific EC2 instance in a region
Analyze the month-over-month cost trend for a particular EC2 instance in a specific region to track its financial impact over time.

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
  aws_cost_by_resource_monthly
where
  resource_id = 'i-1234567890abcdef0'
  and dimension_key = 'REGION'
  and dimension_value = 'us-east-1'
  and period_start >= date_trunc('year', current_date)
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
  aws_cost_by_resource_monthly
where
  resource_id = 'i-1234567890abcdef0'
  and dimension_key = 'REGION'
  and dimension_value = 'us-east-1'
  and period_start >= date('now', 'start of year')
order by
  period_start desc;
```

### Top 10 most expensive resources this month by service
Identify the resources that are generating the highest costs in the current month within a specific service to focus cost optimization efforts.

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
  aws_cost_by_resource_monthly
where
  period_start = date_trunc('month', current_date)
  and dimension_key = 'SERVICE'
  and dimension_value = 'Amazon Elastic Compute Cloud - Compute'
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
  aws_cost_by_resource_monthly
where
  period_start = date('now', 'start of month')
  and dimension_key = 'SERVICE'
  and dimension_value = 'Amazon Elastic Compute Cloud - Compute'
order by
  blended_cost_amount desc
limit 10;
```

### Cost trend analysis by resource type and linked account
Analyze how costs for different types of resources have changed over time within a specific linked account to identify trends and potential areas for optimization.

```sql+postgres
with monthly_costs as (
  select
    split_part(resource_id, '/', 1) as resource_type,
    period_start,
    sum(blended_cost_amount)::numeric::money as total_cost,
    sum(net_unblended_cost_amount)::numeric::money as net_cost,
    sum(usage_quantity_amount) as total_usage,
    sum(normalized_usage_amount) as total_normalized_usage,
    count(distinct resource_id) as resource_count
  from
    aws_cost_by_resource_monthly
  where
    period_start >= date_trunc('year', current_date)
    and dimension_key = 'LINKED_ACCOUNT'
    and dimension_value = '123456789012'
  group by
    resource_type,
    period_start
)
select
  resource_type,
  period_start,
  total_cost,
  net_cost,
  total_usage,
  total_normalized_usage,
  resource_count,
  lag(total_cost) over (partition by resource_type order by period_start) as previous_month_cost
from
  monthly_costs
order by
  resource_type,
  period_start desc;
```

```sql+sqlite
with monthly_costs as (
  select
    substr(resource_id, 1, instr(resource_id, '/') - 1) as resource_type,
    period_start,
    round(sum(cast(blended_cost_amount as decimal)), 2) as total_cost,
    round(sum(cast(net_unblended_cost_amount as decimal)), 2) as net_cost,
    sum(usage_quantity_amount) as total_usage,
    sum(normalized_usage_amount) as total_normalized_usage,
    count(distinct resource_id) as resource_count
  from
    aws_cost_by_resource_monthly
  where
    period_start >= date('now', 'start of year')
    and dimension_key = 'LINKED_ACCOUNT'
    and dimension_value = '123456789012'
  group by
    resource_type,
    period_start
)
select
  resource_type,
  period_start,
  total_cost,
  net_cost,
  total_usage,
  total_normalized_usage,
  resource_count,
  lag(total_cost) over (partition by resource_type order by period_start) as previous_month_cost
from
  monthly_costs
order by
  resource_type,
  period_start desc;
```

### Resources with significant month-over-month cost increases by region
Identify resources that have experienced substantial cost increases compared to the previous month within a specific region, which might indicate potential issues or optimization opportunities.

```sql+postgres
with monthly_costs as (
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
    lag(blended_cost_amount::numeric) over (partition by resource_id order by period_start) as previous_month_cost
  from
    aws_cost_by_resource_monthly
  where
    period_start >= date_trunc('year', current_date)
    and dimension_key = 'REGION'
    and dimension_value = 'us-east-1'
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
  previous_month_cost::money as previous_month_cost,
  ((cost - previous_month_cost) / nullif(previous_month_cost, 0) * 100)::numeric(10,2) as cost_increase_percent
from
  monthly_costs
where
  previous_month_cost > 0
  and cost > previous_month_cost * 1.5
order by
  cost_increase_percent desc;
```

```sql+sqlite
with monthly_costs as (
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
    lag(cast(blended_cost_amount as decimal)) over (partition by resource_id order by period_start) as previous_month_cost
  from
    aws_cost_by_resource_monthly
  where
    period_start >= date('now', 'start of year')
    and dimension_key = 'REGION'
    and dimension_value = 'us-east-1'
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
  previous_month_cost,
  round(((cost - previous_month_cost) / case when previous_month_cost = 0 then null else previous_month_cost end * 100), 2) as cost_increase_percent
from
  monthly_costs
where
  previous_month_cost > 0
  and cost > previous_month_cost * 1.5
order by
  cost_increase_percent desc;
```
