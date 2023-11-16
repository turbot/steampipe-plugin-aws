---
title: "Table: aws_cost_by_service_monthly - Query AWS Cost Explorer Service using SQL"
description: "Allows users to query AWS Cost Explorer Service for monthly cost breakdown by service. This table provides details such as the service name, the cost associated with it, and the currency code."
---

# Table: aws_cost_by_service_monthly - Query AWS Cost Explorer Service using SQL

The `aws_cost_by_service_monthly` table in Steampipe provides information about the monthly cost breakdown by service within AWS Cost Explorer. This table allows financial analysts, DevOps engineers, and other stakeholders to query cost-specific details, including the service name, the cost associated with it, and the currency code. Users can utilize this table to gather insights on cost management, such as tracking AWS expenses, identifying cost trends, and auditing. The schema outlines the various attributes of the cost information, including the service name, cost, and currency code.

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage.  The `aws_cost_by_service_monthly` table provides a simplified view of cost for services in your account (or all linked accounts when run against the organization master), summarized by month, for the last year.  

Note that [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request will incur a cost of $0.01.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cost_by_service_monthly` table, you can use the `.inspect aws_cost_by_service_monthly` command in Steampipe.

### Key columns:

- `service_name`: This is the name of the AWS service. It can be used to join with other tables that contain service-specific information.
- `cost`: This column shows the cost associated with the specific AWS service. It is useful for cost analysis and budgeting.
- `currency_code`: This column indicates the currency in which the cost is expressed. This is important when dealing with costs across different regions and currencies.

## Examples

### Basic info

```sql
select
  service,
  period_start,
  blended_cost_amount::numeric::money,
  unblended_cost_amount::numeric::money,
  amortized_cost_amount::numeric::money,
  net_unblended_cost_amount::numeric::money,
  net_amortized_cost_amount::numeric::money
from 
  aws_cost_by_service_monthly
order by
  service,
  period_start;
```



### Min, Max, and average monthly unblended_cost_amount by service

```sql
select
  service,
  min(unblended_cost_amount)::numeric::money as min,
  max(unblended_cost_amount)::numeric::money as max,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_service_monthly
group by
  service
order by
  service;
```

### Top 10 most expensive service (by average monthly unblended_cost_amount)

```sql
select
  service,
  sum(unblended_cost_amount)::numeric::money as sum,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_service_monthly
group by
  service
order by
  average desc
limit 10;
```


### Top 10 most expensive service (by total monthly unblended_cost_amount)

```sql
select
  service,
  sum(unblended_cost_amount)::numeric::money as sum,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_service_monthly
group by
  service
order by
  sum desc
limit 10;
```


### Ranked - Most expensive month (unblended_cost_amount) by service

```sql
with ranked_costs as (
  select
    service,
    period_start,
    unblended_cost_amount::numeric::money,
    rank() over(partition by service order by unblended_cost_amount desc)
  from 
    aws_cost_by_service_monthly
)
select * from ranked_costs where rank = 1
```



###  Month on month growth (unblended_cost_amount) by service

```sql
with cost_data as (
  select
    service,
    period_start,
    unblended_cost_amount as this_month,
    lag(unblended_cost_amount,1) over(partition by service order by period_start desc) as previous_month
  from 
    aws_cost_by_service_monthly
)
select
    service,
    period_start,
    this_month::numeric::money,
    previous_month::numeric::money,
    case 
      when previous_month = 0 and this_month = 0  then 0
      when previous_month = 0 then 999
      else round((100 * ( (this_month - previous_month) / previous_month))::numeric, 2) 
    end as percent_change
from
  cost_data
order by
  service,
  period_start;
```