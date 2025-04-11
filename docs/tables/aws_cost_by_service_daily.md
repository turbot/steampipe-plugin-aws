---
title: "Steampipe Table: aws_cost_by_service_daily - Query AWS Cost Explorer using SQL"
description: "Allows users to query AWS Cost Explorer to retrieve daily cost breakdown by AWS service."
folder: "Cost Explorer"
---

# Table: aws_cost_by_service_daily - Query AWS Cost Explorer using SQL

The AWS Cost Explorer is a tool that allows you to visualize, understand, and manage your AWS costs and usage over time. It provides data about your cost drivers and usage trends, and enables you to drill down into your cost data to identify specific cost allocation tags or accounts in your organization. You can use it to track your daily AWS costs by service, making it easier to manage your AWS spending.

## Table Usage Guide

The `aws_cost_by_service_daily` table in Steampipe provides you with information about the daily cost breakdown by AWS service within AWS Cost Explorer. This table allows you, as a financial analyst or cloud administrator, to query cost-specific details, including total cost, unit, and service name on a daily basis. You can utilize this table to track your spending on AWS services, monitor cost trends, and identify potential cost-saving opportunities. The schema outlines the various attributes of your cost data, including your linked account, service, currency, and amount.

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage. The `aws_cost_by_service_daily` table provides you with a simplified view of cost for services in your account (or all linked accounts when run against the organization master), summarized by day, for the last year. 

**Important Notes**
- The [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request you make will incur a cost of $0.01.

## Examples

### Basic info
Explore your daily AWS costs by service over a period of time. This query helps you track and analyze your expenditure, aiding in better financial management and budget planning.

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
  aws_cost_by_service_daily
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
  aws_cost_by_service_daily
order by
  service,
  period_start;
```

### Min, Max, and average daily unblended_cost_amount by service
This query is useful for gaining insights into the range and average of daily costs associated with different services in AWS. It can assist in identifying areas of high expenditure and evaluating cost efficiency.

```sql+postgres
select
  service,
  min(unblended_cost_amount)::numeric::money as min,
  max(unblended_cost_amount)::numeric::money as max,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_service_daily
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
  aws_cost_by_service_daily
group by
  service
order by
  service;
```

### Top 10 most expensive service (by average daily unblended_cost_amount)
Discover the segments that are driving your AWS costs by identifying the top 10 most expensive services based on their average daily costs. This helps in managing your budget more effectively and strategically allocating resources.

```sql+postgres
select
  service,
  sum(unblended_cost_amount)::numeric::money as sum,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_service_daily
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
  aws_cost_by_service_daily
group by
  service
order by
  average desc
limit 10;
```

### Top 10 most expensive service (by total daily unblended_cost_amount)
Determine the areas in which your AWS services are costing the most by identifying the top 10 services with the highest daily costs. This can help in optimizing resources and budgeting by focusing on the most expensive services.

```sql+postgres
select
  service,
  sum(unblended_cost_amount)::numeric::money as sum,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_service_daily
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
  aws_cost_by_service_daily
group by
  service
order by
  sum desc
limit 10;
```


### Ranked - Top 10 Most expensive days (unblended_cost_amount) by service
This query is used to identify the top 10 days with the highest expenses for each service. This information can be helpful in managing budgets and identifying potential cost-saving opportunities.

```sql+postgres
with ranked_costs as (
  select
    service,
    period_start,
    unblended_cost_amount::numeric::money,
    rank() over(partition by service order by unblended_cost_amount desc)
  from 
    aws_cost_by_service_daily
)
select * from ranked_costs where rank <= 10;
```

```sql+sqlite
Error: SQLite does not support the rank window function.
```