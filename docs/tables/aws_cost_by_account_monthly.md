---
title: "Steampipe Table: aws_cost_by_account_monthly - Query AWS Cost Explorer Service using SQL"
description: "Allows users to query monthly AWS costs per account. It provides cost details for each AWS account, allowing users to monitor and manage their AWS spending."
folder: "Cost Explorer"
---

# Table: aws_cost_by_account_monthly - Query AWS Cost Explorer Service using SQL

The AWS Cost Explorer Service provides insights into your AWS costs and usage. It enables you to visualize, understand, and manage your AWS costs and usage over time. You can use it to query your monthly AWS costs by account using SQL.

## Table Usage Guide

The `aws_cost_by_account_monthly` table in Steampipe provides you with information about your monthly AWS costs per account. This table allows you, as a financial analyst or DevOps engineer, to query cost-specific details, including the total amount spent, the currency code, and the associated AWS account. You can utilize this table to gain insights on your AWS spending and to manage your budget more effectively. The schema outlines the various attributes of your AWS cost, including the account ID, the month, the total amount, and the currency code.

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage. The `aws_cost_by_account_monthly` table provides a simplified view of cost for your account (or all linked accounts when run against the organization master), summarized by month, for the last year.  

**Important Notes**
- The [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request you make will incur a cost of $0.01.

## Examples

### Basic info
This query allows you to analyze the monthly costs associated with each linked account on AWS. It helps in understanding the financial impact of different accounts and provides insights for better cost management.

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
  aws_cost_by_account_monthly
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
  aws_cost_by_account_monthly
order by
  linked_account_id,
  period_start;
```



### Min, Max, and average monthly unblended_cost_amount by account
Analyze your AWS accounts' monthly expenditure to identify the minimum, maximum, and average costs. This information can help in budgeting and managing your cloud expenses more effectively.

```sql+postgres
select
  linked_account_id,
  min(unblended_cost_amount)::numeric::money as min,
  max(unblended_cost_amount)::numeric::money as max,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_account_monthly
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
  aws_cost_by_account_monthly
group by
  linked_account_id
order by
  linked_account_id;
```


### Ranked - Most expensive months (unblended_cost_amount) by account
Analyze your spending patterns by identifying the months with the highest costs for each linked AWS account. This can help manage your budget by highlighting periods of increased expenditure.

```sql+postgres
select
  linked_account_id,
  period_start,
  unblended_cost_amount::numeric::money,
  rank() over(partition by linked_account_id order by unblended_cost_amount desc)
from 
  aws_cost_by_account_monthly;
```

```sql+sqlite
Error: SQLite does not support the rank window function.
```

### Month on month growth (unblended_cost_amount) by account
This query is designed to analyze monthly expenditure trends across different accounts. It helps users identify any significant changes in costs, which can be useful for budgeting and cost management purposes.

```sql+postgres
with cost_data as (
  select
    linked_account_id,
    period_start,
    unblended_cost_amount as this_month,
    lag(unblended_cost_amount,-1) over(partition by linked_account_id order by period_start desc) as previous_month
  from 
    aws_cost_by_account_monthly
)
select
    linked_account_id,
    period_start,
    this_month::numeric::money,
    previous_month::numeric::money,
    round((100 * ( (this_month - previous_month) / previous_month))::numeric, 2) as percent_change
from
  cost_data
order by
  linked_account_id,
  period_start;
```

```sql+sqlite
with cost_data as (
  select
    linked_account_id,
    period_start,
    unblended_cost_amount as this_month,
    lag(unblended_cost_amount, -1) over(partition by linked_account_id order by period_start desc) as previous_month
  from 
    aws_cost_by_account_monthly
)
select
    linked_account_id,
    period_start,
    this_month,
    previous_month,
    round(100 * (this_month - previous_month) / previous_month, 2) as percent_change
from
  cost_data
order by
  linked_account_id;
```