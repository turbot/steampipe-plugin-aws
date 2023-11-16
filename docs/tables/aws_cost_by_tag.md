---
title: "Table: aws_cost_by_tag - Query AWS Cost Explorer using SQL"
description: "Allows users to query AWS Cost Explorer to obtain cost allocation tags and associated costs."
---

# Table: aws_cost_by_tag - Query AWS Cost Explorer using SQL

The `aws_cost_by_tag` table in Steampipe provides information about cost allocation tags and associated costs within AWS Cost Explorer. This table allows financial analysts, cloud economists, and DevOps engineers to query cost-specific details, including costs associated with each tag. Users can utilize this table to gather insights on cost allocation, such as identifying the most expensive tags, tracking costs of specific projects, departments, or services, and more. The schema outlines the various attributes of the cost allocation tag, including the tag key, cost, and currency.

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage. The `aws_cost_by_tag` table provides a simplified view of cost by tags in your account. One must specify a granularity (`MONTHLY`, `DAILY`) and `tag_key_1` to query the table, however, `tag_key_2` is optional.

Note that [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request will incur a cost of $0.01.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cost_by_tag` table, you can use the `.inspect aws_cost_by_tag` command in Steampipe.

**Key columns**:

- `tag_key`: This is the key of the cost allocation tag. It is important as it allows users to identify which tag the costs are associated with.
- `cost`: This is the cost associated with the tag. It is useful for tracking and managing costs.
- `currency`: This indicates the currency in which the cost is measured. It is important for accurate financial analysis across different regions and currencies.

## Examples

### Basic info

```sql
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

### Min, Max, and average daily unblended_cost_amount by tag

```sql
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

### Ranked - Top 10 Most expensive days (unblended_cost_amount) by tag

```sql
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
