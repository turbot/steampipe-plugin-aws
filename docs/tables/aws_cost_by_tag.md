---
title: "Steampipe Table: aws_cost_by_tag - Query AWS Cost Explorer using SQL"
description: "Allows users to query AWS Cost Explorer to obtain cost allocation tags and associated costs."
folder: "Cost Explorer"
---

# Table: aws_cost_by_tag - Query AWS Cost Explorer using SQL

The AWS Cost Explorer is a tool that enables you to view and analyze your costs and usage. You can explore your AWS costs using an interface that allows you to break down costs by AWS service, linked account, tag, and many other dimensions. Through the AWS Cost Explorer API, you can directly access this data and use it to create your own cost management applications.

## Table Usage Guide

The `aws_cost_by_tag` table in Steampipe provides you with information about cost allocation tags and associated costs within AWS Cost Explorer. This table allows you, as a financial analyst, cloud economist, or DevOps engineer, to query cost-specific details, including costs associated with each tag. You can utilize this table to gather insights on cost allocation, such as identifying the most expensive tags, tracking costs of specific projects, departments, or services, and more. The schema outlines the various attributes of the cost allocation tag, including the tag key, cost, and currency.

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage. The `aws_cost_by_tag` table provides you with a simplified view of cost by tags in your account. You must specify a granularity (`MONTHLY`, `DAILY`) and `tag_key_1` to query the table, however, `tag_key_2` is optional.

**Important Notes**

- The [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request you make will incur a cost of $0.01.

## Examples

### Basic info
This query is used to gain insights into the daily cost breakdown of AWS services, based on specific tags. It is particularly useful for tracking and managing costs, especially in scenarios where resources are tagged by project, department, or any other category for cost allocation.

```sql+postgres
select
  tag_key_1,
  tag_value_1,
  period_start,
  blended_cost_amount::numeric::money,
  unblended_cost_amount::numeric::money,
  amortized_cost_amount::numeric::money,
  net_unblended_cost_amount::numeric::money,
  net_amortized_cost_amount::numeric::money
from
  aws_cost_by_tag
where
  granularity = 'DAILY'
and
  tag_key_1 = 'Name';
```

```sql+sqlite
select
  tag_key_1,
  tag_value_1,
  period_start,
  CAST(blended_cost_amount AS NUMERIC) AS blended_cost_amount,
  CAST(unblended_cost_amount AS NUMERIC) AS unblended_cost_amount,
  CAST(amortized_cost_amount AS NUMERIC) AS amortized_cost_amount,
  CAST(net_unblended_cost_amount AS NUMERIC) AS net_unblended_cost_amount,
  CAST(net_amortized_cost_amount AS NUMERIC) AS net_amortized_cost_amount
from
  aws_cost_by_tag
where
  granularity = 'DAILY'
and
  tag_key_1 = 'Name';
```

### Min, Max, and average daily unblended_cost_amount by tag
Discover the segments that have the lowest, highest, and average daily costs associated with a specific tag. This is useful for tracking and managing AWS costs on a day-to-day basis by identifying areas where spending is concentrated.

```sql+postgres
select
  tag_key_1,
  tag_value_1,
  min(unblended_cost_amount)::numeric::money as min,
  max(unblended_cost_amount)::numeric::money as max,
  avg(unblended_cost_amount)::numeric::money as average
from
  aws_cost_by_tag
where
  granularity = 'DAILY'
and
  tag_key_1 = 'Name'
group by
  tag_key_1, tag_value_1;
```

```sql+sqlite
select
  tag_key_1,
  tag_value_1,
  min(unblended_cost_amount) as min,
  max(unblended_cost_amount) as max,
  avg(unblended_cost_amount) as average
from
  aws_cost_by_tag
where
  granularity = 'DAILY'
and
  tag_key_1 = 'Name'
group by
  tag_key_1, tag_value_1;
```

### Ranked - Top 10 Most expensive days (unblended_cost_amount) by tag
Discover the segments that are the top 10 most costly, based on daily expenditures, to identify potential areas of cost reduction. This is particularly useful for those looking to optimize their resource utilization and manage their budget effectively.

```sql+postgres
with ranked_costs as
(
  select
    tag_key_1,
    tag_value_1,
    period_start,
    unblended_cost_amount::numeric::money,
    rank() over(partition by tag_key_1
  order by
    unblended_cost_amount desc)
  from
    aws_cost_by_tag
  where
    granularity = 'DAILY'
    and tag_key_1 = 'Name'
)
select
  *
from
  ranked_costs
where
  rank <= 10;
```

```sql+sqlite
Error: SQLite does not support the rank window function.
```