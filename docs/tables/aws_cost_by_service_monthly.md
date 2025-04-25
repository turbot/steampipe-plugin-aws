---
title: "Steampipe Table: aws_cost_by_service_monthly - Query AWS Cost Explorer Service using SQL"
description: "Allows users to query AWS Cost Explorer Service for monthly cost breakdown by service. This table provides details such as the service name, the cost associated with it, and the currency code."
folder: "Cost Explorer"
---

# Table: aws_cost_by_service_monthly - Query AWS Cost Explorer Service using SQL

The AWS Cost Explorer Service provides detailed information about your AWS costs, enabling you to analyze your costs and usage over time. You can use it to identify trends, isolate cost drivers, and detect anomalies. With SQL queries, you can retrieve monthly cost data specific to each AWS service.

## Table Usage Guide

The `aws_cost_by_service_monthly` table in Steampipe provides you with information about the monthly cost breakdown by service within AWS Cost Explorer. This table allows you, as a financial analyst, DevOps engineer, or other stakeholder, to query cost-specific details, including the service name, the cost associated with it, and the currency code. You can utilize this table to gather insights on cost management, such as tracking AWS expenses, identifying cost trends, and auditing. The schema outlines the various attributes of the cost information, including the service name, cost, and currency code.

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage.  The `aws_cost_by_service_monthly` table provides you with a simplified view of cost for services in your account (or all linked accounts when run against the organization master), summarized by month, for the last year.  

**Important Notes**

- The [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request you make will incur a cost of $0.01.

## Examples

### Basic info
Explore which AWS services have the highest costs over time. This query is useful in identifying potential areas for cost reduction through service optimization or consolidation.

```sql+postgres
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

```sql+sqlite
select
  service,
  period_start,
  cast(blended_cost_amount as decimal),
  cast(unblended_cost_amount as decimal),
  cast(amortized_cost_amount as decimal),
  cast(net_unblended_cost_amount as decimal),
  cast(net_amortized_cost_amount as decimal)
from 
  aws_cost_by_service_monthly
order by
  service,
  period_start;
```



### Min, Max, and average monthly unblended_cost_amount by service
Explore which AWS services have the lowest, highest, and average monthly costs, providing a clear understanding of your AWS expenditure. This can help in budgeting and identifying services that may be costing more than expected.

```sql+postgres
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

```sql+sqlite
select
  service,
  min(unblended_cost_amount) as min,
  max(unblended_cost_amount) as max,
  avg(unblended_cost_amount) as average
from 
  aws_cost_by_service_monthly
group by
  service
order by
  service;
```

### Top 10 most expensive service (by average monthly unblended_cost_amount)
Discover the segments that are incurring the highest average monthly costs on your AWS account. This information can be crucial for budgeting and cost management strategies, helping you to identify areas where expenses can be reduced.

```sql+postgres
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

```sql+sqlite
select
  service,
  sum(unblended_cost_amount) as sum,
  avg(unblended_cost_amount) as average
from 
  aws_cost_by_service_monthly
group by
  service
order by
  average desc
limit 10;
```


### Top 10 most expensive service (by total monthly unblended_cost_amount)
This query helps to pinpoint the top 10 most costly services in terms of total monthly unblended cost. It is useful for gaining insights into where the majority of your AWS costs are coming from, allowing for more informed budgeting and cost management decisions.

```sql+postgres
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

```sql+sqlite
select
  service,
  sum(unblended_cost_amount) as sum,
  avg(unblended_cost_amount) as average
from 
  aws_cost_by_service_monthly
group by
  service
order by
  sum desc
limit 10;
```


### Ranked - Most expensive month (unblended_cost_amount) by service
This query is designed to identify the most costly month for each service in terms of unblended costs. It can be useful for budgeting and cost management, helping to highlight areas where expenses may be unexpectedly high.

```sql+postgres
with ranked_costs as (
  select
    service,
    period_start,
    unblended_cost_amount::numeric::money,
    rank() over(partition by service order by unblended_cost_amount desc)
  from 
    aws_cost_by_service_monthly
)
select * from ranked_costs where rank = 1;
```

```sql+sqlite
Error: SQLite does not support the rank window function.
```

###  Month on month growth (unblended_cost_amount) by service
Analyze your AWS monthly costs to understand the percentage change in expenditure for each service. This could be useful for identifying trends, managing budgets, and making strategic decisions about resource allocation.

```sql+postgres
with cost_data as (
  select
    service,
    period_start,
    unblended_cost_amount as this_month,
    lag(unblended_cost_amount,-1) over(partition by service order by period_start desc) as previous_month
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

```sql+sqlite
with cost_data as (
  select
    service,
    period_start,
    unblended_cost_amount as this_month,
    lag(unblended_cost_amount,-1) over(partition by service order by period_start desc) as previous_month
  from 
    aws_cost_by_service_monthly
)
select
    service,
    period_start,
    this_month,
    previous_month,
    case 
      when previous_month = 0 and this_month = 0  then 0
      when previous_month = 0 then 999
      else round((100 * ( (this_month - previous_month) / previous_month)), 2) 
    end as percent_change
from
  cost_data
order by
  service,
  period_start;
```