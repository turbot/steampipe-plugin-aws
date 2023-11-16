---
title: "Table: aws_cost_by_service_daily - Query AWS Cost Explorer using SQL"
description: "Allows users to query AWS Cost Explorer to retrieve daily cost breakdown by AWS service."
---

# Table: aws_cost_by_service_daily - Query AWS Cost Explorer using SQL

The `aws_cost_by_service_daily` table in Steampipe provides information about daily cost breakdown by AWS service within AWS Cost Explorer. This table allows financial analysts and cloud administrators to query cost-specific details, including total cost, unit, and service name on a daily basis. Users can utilize this table to track spending on AWS services, monitor cost trends, and identify potential cost-saving opportunities. The schema outlines the various attributes of the cost data, including the linked account, service, currency, and amount.

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage.  The `aws_cost_by_service_daily` table provides a simplified view of cost for services in your account (or all linked accounts when run against the organization master), summarized by day, for the last year.  

Note that [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request will incur a cost of $0.01.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cost_by_service_daily` table, you can use the `.inspect aws_cost_by_service_daily` command in Steampipe.

### Key columns:

- `linked_account`: This is the ID of the AWS account. This column is useful for identifying which AWS account the costs are associated with, especially in multi-account AWS environments.
- `service`: This column contains the name of the AWS service. It is useful for identifying which AWS service the costs are associated with.
- `date`: This column contains the date for the daily cost data. It is useful for tracking cost trends over time.

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
  aws_cost_by_service_daily
order by
  service,
  period_start;
```



### Min, Max, and average daily unblended_cost_amount by service

```sql
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

### Top 10 most expensive service (by average daily unblended_cost_amount)

```sql
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


### Top 10 most expensive service (by total daily unblended_cost_amount)

```sql
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


### Ranked - Top 10 Most expensive days (unblended_cost_amount) by service

```sql
with ranked_costs as (
  select
    service,
    period_start,
    unblended_cost_amount::numeric::money,
    rank() over(partition by service order by unblended_cost_amount desc)
  from 
    aws_cost_by_service_daily
)
select * from ranked_costs where rank <= 10
```
