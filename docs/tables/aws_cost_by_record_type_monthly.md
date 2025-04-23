---
title: "Steampipe Table: aws_cost_by_record_type_monthly - Query AWS Cost and Usage Report Records using SQL"
description: "Allows users to query AWS Cost and Usage Report Records on a monthly basis."
folder: "Cost Explorer"
---

# Table: aws_cost_by_record_type_monthly - Query AWS Cost and Usage Report Records using SQL

The AWS Cost and Usage Report service provides comprehensive cost and usage data about your AWS resources, enabling you to manage your costs and optimize your AWS spend. It records the AWS usage data for your accounts and delivers the log files to a specified Amazon S3 bucket. You can query these records using SQL to gain insights into your resource usage and cost.

## Table Usage Guide

The `aws_cost_by_record_type_monthly` table in Steampipe provides you with information about AWS Cost and Usage Report Records, specifically detailing costs incurred by different record types on a monthly basis. This table allows you, whether you're a DevOps engineer or a financial analyst, to query cost-specific details, including service usage, cost allocation, and associated metadata. You can utilize this table to gather insights on AWS costs, such as costs associated with specific AWS services, cost trends over time, and cost allocation across different record types. The schema outlines the various attributes of the cost and usage report record, including the record type, usage type, operation, and cost.

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage.  The `aws_cost_by_record_type_monthly` table provides a simplified view of cost for your account (or all linked accounts when run against the organization master) as per record types (fees, usage, costs, tax refunds, and credits), summarized by month, for the last year.  

**Important Notes**
- The [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request will incur a cost of $0.01 for you.

## Examples

### Basic info
Gain insights into your AWS cost trends by analyzing monthly expenses. This query helps in understanding the cost incurred over time, aiding in effective budget planning and cost management.

```sql+postgres
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
  aws_cost_by_record_type_monthly
order by
  linked_account_id,
  period_start;
```

```sql+sqlite
select
  linked_account_id,
  record_type,
  period_start,
  CAST(blended_cost_amount AS REAL) AS blended_cost_amount,
  CAST(unblended_cost_amount AS REAL) AS unblended_cost_amount,
  CAST(amortized_cost_amount AS REAL) AS amortized_cost_amount,
  CAST(net_unblended_cost_amount AS REAL) AS net_unblended_cost_amount,
  CAST(net_amortized_cost_amount AS REAL) AS net_amortized_cost_amount
from 
  aws_cost_by_record_type_monthly
order by
  linked_account_id,
  period_start;
```

### Min, Max, and average monthly unblended_cost_amount by account and record type
Explore which linked accounts have the highest, lowest, and average monthly costs, grouped by record type. This can help in understanding the cost distribution and identifying any unusual spending patterns.

```sql+postgres
select
  linked_account_id,
  record_type,
  min(unblended_cost_amount)::numeric::money as min,
  max(unblended_cost_amount)::numeric::money as max,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_record_type_monthly
group by
  linked_account_id,
  record_type
order by
  linked_account_id;
```

```sql+sqlite
select
  linked_account_id,
  record_type,
  min(unblended_cost_amount) as min,
  max(unblended_cost_amount) as max,
  avg(unblended_cost_amount) as average
from 
  aws_cost_by_record_type_monthly
group by
  linked_account_id,
  record_type
order by
  linked_account_id;
```

### Ranked - Most expensive months (unblended_cost_amount) by account and record type
Explore which months have been the most costly for each account and record type. This can aid in identifying trends and planning future budgeting strategies.

```sql+postgres
select
  linked_account_id,
  record_type,
  period_start,
  unblended_cost_amount::numeric::money,
  rank() over(partition by linked_account_id, record_type order by unblended_cost_amount desc)
from 
  aws_cost_by_record_type_monthly;
```

```sql+sqlite
select
  linked_account_id,
  record_type,
  period_start,
  unblended_cost_amount,
  (
    select count(*) + 1 
    from aws_cost_by_record_type_monthly as b
    where 
      a.linked_account_id = b.linked_account_id and 
      a.record_type = b.record_type and 
      a.unblended_cost_amount < b.unblended_cost_amount
  )
from 
  aws_cost_by_record_type_monthly as a;
```