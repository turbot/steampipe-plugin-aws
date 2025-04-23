---
title: "Steampipe Table: aws_cost_by_account_daily - Query AWS Cost Explorer using SQL"
description: "Allows users to query daily AWS costs by account. This table provides an overview of AWS usage and cost data for each AWS account on a daily basis."
folder: "Cost Explorer"
---

# Table: aws_cost_by_account_daily - Query AWS Cost Explorer using SQL

The AWS Cost Explorer is a service that allows you to visualize, understand, and manage your AWS costs and usage over time. It provides detailed information about your costs and usage, including both AWS service usage and the costs associated with your usage. You can use Cost Explorer to identify trends, pinpoint cost drivers, and detect anomalies.

## Table Usage Guide

The `aws_cost_by_account_daily` table in Steampipe provides you with information about your daily AWS costs for each of your accounts within AWS Cost Explorer. This table allows you, as a financial analyst, cloud economist, or DevOps engineer, to query daily cost-specific details, including cost usage, unblended costs, and associated metadata. You can utilize this table to gather insights on your daily AWS spending, such as cost trends, cost spikes, and cost predictions. The schema outlines the various attributes of your daily cost, including your linked account, service, currency code, and cost usage details.

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage. The `aws_cost_by_account_daily` table provides you with a simplified view of cost for your account (or all linked accounts when run against the organization master), summarized by day, for the last year.

**Important Notes**
- The [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request you make will incur a cost of $0.01.

## Examples

### Basic info
This example allows users to gain insights into their daily AWS cost by account. It's useful for tracking and analyzing cost trends over time, helping to manage and optimize cloud spending.

```sql+postgres
select
  linked_account_id,
  period_start,
  blended_cost_amount::numeric::money,
  unblended_cost_amount::numeric::money,
  amortized_cost_amount::numeric::money,
  net_unblended_cost_amount::numeric::money,
  net_amortized_cost_amount::numeric::money
from 
  aws_cost_by_account_daily
order by
  linked_account_id,
  period_start;
```

```sql+sqlite
select
  linked_account_id,
  period_start,
  CAST(blended_cost_amount AS REAL) AS blended_cost_amount,
  CAST(unblended_cost_amount AS REAL) AS unblended_cost_amount,
  CAST(amortized_cost_amount AS REAL) AS amortized_cost_amount,
  CAST(net_unblended_cost_amount AS REAL) AS net_unblended_cost_amount,
  CAST(net_amortized_cost_amount AS REAL) AS net_amortized_cost_amount
from 
  aws_cost_by_account_daily
order by
  linked_account_id,
  period_start;
```

### Min, Max, and average daily unblended_cost_amount by account
Analyze your AWS accounts to understand the minimum, maximum, and average daily costs. This is useful for monitoring the financial performance of different accounts and identifying potential areas for cost optimization.

```sql+postgres
select
  linked_account_id,
  min(unblended_cost_amount)::numeric::money as min,
  max(unblended_cost_amount)::numeric::money as max,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_account_daily
group by
  linked_account_id
order by
  linked_account_id;
```

```sql+sqlite
select
  linked_account_id,
  min(unblended_cost_amount) as min,
  max(unblended_cost_amount) as max,
  avg(unblended_cost_amount) as average
from 
  aws_cost_by_account_daily
group by
  linked_account_id
order by
  linked_account_id;
```


### Ranked - Top 10 Most expensive days (unblended_cost_amount) by account
Explore the days where the cost was at its highest for each account. This query is useful for identifying potential anomalies or trends in spending, enabling more effective financial management.

```sql+postgres
with ranked_costs as (
  select
    linked_account_id,
    period_start,
    unblended_cost_amount::numeric::money,
    rank() over(partition by linked_account_id order by unblended_cost_amount desc)
  from 
    aws_cost_by_account_daily
)
select * from ranked_costs where rank <= 10;
```

```sql+sqlite
Error: SQLite does not support the rank window function.
```