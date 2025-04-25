---
title: "Steampipe Table: aws_cost_by_record_type_daily - Query AWS Cost and Usage Report using SQL"
description: "Allows users to query daily AWS cost data by record type. This table provides information about AWS costs incurred per record type on a daily basis."
folder: "Cost Explorer"
---

# Table: aws_cost_by_record_type_daily - Query AWS Cost and Usage Report using SQL

The AWS Cost and Usage Report is a comprehensive resource that provides detailed information about your AWS costs. It allows you to view your AWS usage and costs for each service category used by your accounts and by specific cost allocation tags. By querying this report, you can gain insights into your AWS spending and optimize your resource utilization.

## Table Usage Guide

The `aws_cost_by_record_type_daily` table in Steampipe provides you with information about AWS costs incurred per record type on a daily basis. This table allows you as a financial analyst, DevOps engineer, or other professional to query cost-specific details, including the linked account, service, usage type, and operation. You can utilize this table to gather insights on cost distribution, such as costs associated with different services, usage types, and operations. The schema outlines the various attributes of the cost record, including the record id, record type, billing period start date, and cost.

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage.  The `aws_cost_by_record_type_daily` table provides a simplified view of cost for your account (or all linked accounts when run against the organization master) as per record types (fees, usage, costs, tax refunds, and credits), summarized by day, for the last year.  

**Important Notes**
- The [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request you make will incur a cost of $0.01.

## Examples

### Basic info
Determine the areas in which your AWS account incurs costs on a daily basis. This query helps you understand your spending patterns by breaking down costs into different categories, allowing you to manage your AWS resources more efficiently.

```sql+postgres
select
  linked_account_id,
  record_type,
  period_start,
  blended_cost_amount::numeric::money,
  unblended_cost_amount::numeric::money,
  amortized_cost_amount::numeric::money,
  net_unblended_cost_amount::numeric::money,
  net_amortized_cost_amount::numeric::money
from 
  aws_cost_by_record_type_daily
order by
  linked_account_id,
  period_start;
```

```sql+sqlite
select
  linked_account_id,
  record_type,
  period_start,
  CAST(blended_cost_amount AS REAL) AS blended_cost_amount,
  CAST(unblended_cost_amount AS REAL) AS unblended_cost_amount,
  CAST(amortized_cost_amount AS REAL) AS amortized_cost_amount,
  CAST(net_unblended_cost_amount AS REAL) AS net_unblended_cost_amount,
  CAST(net_amortized_cost_amount AS REAL) AS net_amortized_cost_amount
from 
  aws_cost_by_record_type_daily
order by
  linked_account_id,
  period_start;
```

### Min, Max, and average daily unblended_cost_amount by account and record type
Determine the areas in which you have minimum, maximum, and average daily costs associated with different accounts and record types. This can help you identify potential cost-saving opportunities and better manage your resources.

```sql+postgres
select
  linked_account_id,
  record_type,
  min(unblended_cost_amount)::numeric::money as min,
  max(unblended_cost_amount)::numeric::money as max,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_record_type_daily
group by
  linked_account_id,
  record_type
order by
  linked_account_id;
```

```sql+sqlite
select
  linked_account_id,
  record_type,
  min(unblended_cost_amount) as min,
  max(unblended_cost_amount) as max,
  avg(unblended_cost_amount) as average
from 
  aws_cost_by_record_type_daily
group by
  linked_account_id,
  record_type
order by
  linked_account_id;
```

### Ranked - Top 10 Most expensive days (unblended_cost_amount) by account and record type
Determine the days with the highest expenses, grouped by account and record type. This query can help in cost optimization by identifying the top 10 most expensive days, allowing for better budget management and resource allocation.

```sql+postgres
with ranked_costs as (
  select
    linked_account_id,
    record_type,
    period_start,
    unblended_cost_amount::numeric::money,
    rank() over(partition by linked_account_id, record_type order by unblended_cost_amount desc)
  from 
    aws_cost_by_record_type_daily
)
select * from ranked_costs where rank <= 10;
```

```sql+sqlite
Error: SQLite does not support the rank window function.
```