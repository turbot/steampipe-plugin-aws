---
title: "Table: aws_cost_usage - Query AWS Cost Explorer Service Cost and Usage using SQL"
description: "Allows users to query Cost and Usage data from AWS Cost Explorer Service to monitor, track, and manage AWS costs and usage over time."
---

# Table: aws_cost_usage - Query AWS Cost Explorer Service Cost and Usage using SQL

The `aws_cost_usage` table in Steampipe provides information about cost and usage data from AWS Cost Explorer Service. This table allows financial analysts and cloud architects to query cost and usage details, including cost allocation tags, service usage, cost usage, and associated metadata. Users can utilize this table to gather insights on cost and usage, such as cost per service, usage per service, verification of cost allocation tags, and more. The schema outlines the various attributes of the cost and usage data, including the time period, unblended cost, usage type, and associated tags.

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage.  The `aws_cost_usage` table provides a simplified yet flexible view of cost for your account (or all linked accounts when run against the organization master).  You must specify a granularity (`MONTHLY`, `DAILY`), and 2 dimension types (`AZ` , `INSTANCE_TYPE`, `LEGAL_ENTITY_NAME`, `LINKED_ACCOUNT`, `OPERATION`, `PLATFORM`, `PURCHASE_TYPE`, `SERVICE`, `TENANCY`, `RECORD_TYPE`, and `USAGE_TYPE`)

This tables requires an '=' qualifier for all of the following columns: granularity,dimension_type_1,dimension_type_2

Note that [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request will incur a cost of $0.01.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cost_usage` table, you can use the `.inspect aws_cost_usage` command in Steampipe.

### Key columns:

- `identity_line_item_operation`: This column indicates the operation that the line item covers, such as 'RunInstances' for EC2. It is useful for tracking the specific AWS operation that incurs costs.

- `product_product_name`: This column indicates the name of the AWS product that the line item applies to, such as Amazon Simple Storage Service or Amazon EC2. It is important for understanding which AWS services are incurring costs.

- `line_item_usage_type`: This column indicates the usage type for a specific line item. It helps in understanding the type of usage (like data transfer, instance hours) that incurs costs.

## Examples

### Monthly net unblended cost by account and service

```sql
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

### Top 5 most expensive services (net unblended cost) in each account

```sql
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

### Monthly net unblended cost by account and record type

```sql
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

### List monthly discounts and credits by account

```sql
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
