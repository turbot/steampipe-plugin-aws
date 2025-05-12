---
title: "Steampipe Table: aws_cost_usage - Query AWS Cost Explorer Service Cost and Usage using SQL"
description: "Allows users to query Cost and Usage data from AWS Cost Explorer Service to monitor, track, and manage AWS costs and usage over time."
folder: "Cost Explorer"
---

# Table: aws_cost_usage - Query AWS Cost Explorer Service Cost and Usage using SQL

The AWS Cost Explorer Service is a tool that allows you to visualize, understand, and manage your AWS costs and usage over time. It provides detailed information about your costs and usage, including trends, cost drivers, and anomalies. With Cost Explorer, you can filter views by various dimensions such as service, linked account, and tags, and view data for up to the last 13 months.

## Table Usage Guide

The `aws_cost_usage` table in Steampipe provides you with information about cost and usage data from AWS Cost Explorer Service. This table enables you as a financial analyst or cloud architect to query cost and usage details, including cost allocation tags, service usage, cost usage, and associated metadata. You can utilize this table to gather insights on cost and usage, such as cost per service, usage per service, verification of cost allocation tags, and more. The schema outlines the various attributes of the cost and usage data for you, including the time period, unblended cost, usage type, and associated tags.

Amazon Cost Explorer assists you in visualizing, understanding, and managing your AWS costs and usage. The `aws_cost_usage` table offers you a simplified yet flexible view of cost for your account (or all linked accounts when run against the organization master). You need to specify a granularity (`MONTHLY`, `DAILY`), and 2 dimension types (`AZ` , `INSTANCE_TYPE`, `LEGAL_ENTITY_NAME`, `LINKED_ACCOUNT`, `OPERATION`, `PLATFORM`, `PURCHASE_TYPE`, `SERVICE`, `TENANCY`, `RECORD_TYPE`, and `USAGE_TYPE`)

**Important Notes**
- This table requires an '=' qualifier for all of the following columns: granularity, dimension_type_1, dimension_type_2.
- The [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request will incur a cost of $0.01 for you.

## Examples

### Monthly net unblended cost by account and service
Explore the monthly expenditure for each linked account and service in your AWS environment. This query can help you understand your cost trends and identify areas for potential savings.

```sql+postgres
select
  period_start,
  dimension_1 as account_id,
  dimension_2 as service_name,
  net_unblended_cost_amount::numeric::money
from
  aws_cost_usage
where
  granularity = 'MONTHLY'
  and dimension_type_1 = 'LINKED_ACCOUNT'
  and dimension_type_2 = 'SERVICE'
order by
  dimension_1,
  period_start;
```

```sql+sqlite
select
  period_start,
  dimension_1 as account_id,
  dimension_2 as service_name,
  cast(net_unblended_cost_amount as real) as net_unblended_cost_amount
from
  aws_cost_usage
where
  granularity = 'MONTHLY'
  and dimension_type_1 = 'LINKED_ACCOUNT'
  and dimension_type_2 = 'SERVICE'
order by
  dimension_1,
  period_start;
```

### Top 5 most expensive services (net unblended cost) in each account
Identify the top five most costly services in each account to manage and optimize your AWS expenses effectively.

```sql+postgres
with ranked_costs as (
  select
    dimension_1 as account_id,
    dimension_2 as service_name,
    sum(net_unblended_cost_amount)::numeric::money as net_unblended_cost,
    rank() over(partition by dimension_1 order by sum(net_unblended_cost_amount) desc)
  from
    aws_cost_usage
  where
    granularity = 'MONTHLY'
    and dimension_type_1 = 'LINKED_ACCOUNT'
    and dimension_type_2 = 'SERVICE'
  group by
    dimension_1,
    dimension_2
  order by
    dimension_1,
    net_unblended_cost desc
)
select * from ranked_costs where rank <=5
```

```sql+sqlite
Error: SQLite does not support rank window functions.
```

### Monthly net unblended cost by account and record type
Analyze your monthly AWS account costs by record type to better understand your expenses. This can help you identify areas where costs may be reduced or controlled.

```sql+postgres
select
  period_start,
  dimension_1 as account_id,
  dimension_2 as record_type,
  net_unblended_cost_amount::numeric::money
from
  aws_cost_usage
where
  granularity = 'MONTHLY'
  and dimension_type_1 = 'LINKED_ACCOUNT'
  and dimension_type_2 = 'RECORD_TYPE'
order by
  dimension_1,
  period_start;
```

```sql+sqlite
select
  period_start,
  dimension_1 as account_id,
  dimension_2 as record_type,
  CAST(net_unblended_cost_amount AS REAL) AS net_unblended_cost_amount
from
  aws_cost_usage
where
  granularity = 'MONTHLY'
  and dimension_type_1 = 'LINKED_ACCOUNT'
  and dimension_type_2 = 'RECORD_TYPE'
order by
  dimension_1,
  period_start;
```

### List monthly discounts and credits by account
This query allows users to monitor their AWS account's monthly spending by tracking discounts and credits. It's beneficial for budgeting purposes and helps in optimizing cost management strategies.

```sql+postgres
select
  period_start,
  dimension_1 as account_id,
  dimension_2 as record_type,
  net_unblended_cost_amount::numeric::money
from
  aws_cost_usage
where
  granularity = 'MONTHLY'
  and dimension_type_1 = 'LINKED_ACCOUNT'
  and dimension_type_2 = 'RECORD_TYPE'
  and dimension_2 in ('DiscountedUsage', 'Credit')
order by
  dimension_1,
  period_start;
```

```sql+sqlite
select
  period_start,
  dimension_1 as account_id,
  dimension_2 as record_type,
  CAST(net_unblended_cost_amount AS REAL) as net_unblended_cost_amount
from
  aws_cost_usage
where
  granularity = 'MONTHLY'
  and dimension_type_1 = 'LINKED_ACCOUNT'
  and dimension_type_2 = 'RECORD_TYPE'
  and dimension_2 in ('DiscountedUsage', 'Credit')
order by
  dimension_1,
  period_start;
```