---
title: "Table: aws_cost_by_record_type_daily - Query AWS Cost and Usage Report using SQL"
description: "Allows users to query daily AWS cost data by record type. This table provides information about AWS costs incurred per record type on a daily basis."
---

# Table: aws_cost_by_record_type_daily - Query AWS Cost and Usage Report using SQL

The `aws_cost_by_record_type_daily` table in Steampipe provides information about AWS costs incurred per record type on a daily basis. This table allows financial analysts, DevOps engineers, and other professionals to query cost-specific details, including the linked account, service, usage type, and operation. Users can utilize this table to gather insights on cost distribution, such as costs associated with different services, usage types, and operations. The schema outlines the various attributes of the cost record, including the record id, record type, billing period start date, and cost.

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage.  The `aws_cost_by_record_type_daily` table provides a simplified view of cost for your account (or all linked accounts when run against the organization master) as per record types (fees, usage, costs, tax refunds, and credits), summarized by day, for the last year.  

Note that [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request will incur a cost of $0.01.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cost_by_record_type_daily` table, you can use the `.inspect aws_cost_by_record_type_daily` command in Steampipe.

Key columns:

- `record_id`: This is the unique identifier for each cost record. It is useful for tracking individual cost records and can be used to join this table with others that contain cost record IDs.
- `linked_account`: This column contains the ID of the AWS account linked to the cost record. It is crucial for cost allocation and can be used to join this table with others that contain linked account IDs.
- `service`: This column identifies the AWS service associated with the cost record. It is helpful for understanding cost distribution across different AWS services and can be used to join this table with others that contain service names.

## Examples

### Basic info

```sql
select
  linked_account_id,
  record_type,
  period_start,
  blended_cost_amount::numeric::money,
  unblended_cost_amount::numeric::money,
  amortized_cost_amount::numeric::money,
  net_unblended_cost_amount::numeric::money,
  net_amortized_cost_amount::numeric::money
from 
  aws_cost_by_record_type_daily
order by
  linked_account_id,
  period_start;
```

### Min, Max, and average daily unblended_cost_amount by account and record type

```sql
select
  linked_account_id,
  record_type,
  min(unblended_cost_amount)::numeric::money as min,
  max(unblended_cost_amount)::numeric::money as max,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_record_type_daily
group by
  linked_account_id,
  record_type
order by
  linked_account_id;
```

### Ranked - Top 10 Most expensive days (unblended_cost_amount) by account and record type

```sql
with ranked_costs as (
  select
    linked_account_id,
    record_type,
    period_start,
    unblended_cost_amount::numeric::money,
    rank() over(partition by linked_account_id, record_type order by unblended_cost_amount desc)
  from 
    aws_cost_by_record_type_daily
)
select * from ranked_costs where rank <= 10
```
