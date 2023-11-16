---
title: "Table: aws_cost_by_account_daily - Query AWS Cost Explorer using SQL"
description: "Allows users to query daily AWS costs by account. This table provides an overview of AWS usage and cost data for each AWS account on a daily basis."
---

# Table: aws_cost_by_account_daily - Query AWS Cost Explorer using SQL

The `aws_cost_by_account_daily` table in Steampipe provides information about daily AWS costs for each account within AWS Cost Explorer. This table allows financial analysts, cloud economists, and DevOps engineers to query daily cost-specific details, including cost usage, unblended costs, and associated metadata. Users can utilize this table to gather insights on daily AWS spending, such as cost trends, cost spikes, and cost predictions. The schema outlines the various attributes of the daily cost, including the linked account, service, currency code, and cost usage details.

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage.  The `aws_cost_by_account_daily` table provides a simplified view of cost for your account (or all linked accounts when run against the organization master), summarized by day, for the last year.  

Note that [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request will incur a cost of $0.01.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cost_by_account_daily` table, you can use the `.inspect aws_cost_by_account_daily` command in Steampipe.

### Key columns:

- `linked_account`: This column stores the ID of the linked AWS account. It can be used to join this table with other tables that contain account-specific information.
- `service`: This column stores the AWS service name. It can be used to join this table with other tables that contain service-specific information.
- `date`: This column stores the date for the daily cost record. It can be used to join this table with other tables that contain date-specific information.

## Examples

### Basic info

```sql
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



### Min, Max, and average daily unblended_cost_amount by account

```sql
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


### Ranked - Top 10 Most expensive days (unblended_cost_amount) by account

```sql
with ranked_costs as (
  select
    linked_account_id,
    period_start,
    unblended_cost_amount::numeric::money,
    rank() over(partition by linked_account_id order by unblended_cost_amount desc)
  from 
    aws_cost_by_account_daily
)
select * from ranked_costs where rank <= 10
```
