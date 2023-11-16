---
title: "Table: aws_cost_by_service_usage_type_daily - Query AWS Cost Explorer Service usage type daily using SQL"
description: "Allows users to query AWS Cost Explorer Service daily usage type to fetch detailed data about AWS service usage and costs."
---

# Table: aws_cost_by_service_usage_type_daily - Query AWS Cost Explorer Service usage type daily using SQL

The `aws_cost_by_service_usage_type_daily` table in Steampipe provides information about daily usage type and costs for each AWS service within AWS Cost Explorer. This table allows DevOps engineers, financial analysts, and cloud architects to query daily-specific details, including usage amount, usage unit, and the corresponding service cost. Users can utilize this table to gather insights on daily usage and costs, such as identifying high-cost services, tracking usage patterns, and managing AWS expenses. The schema outlines the various attributes of the AWS service cost, including the service name, usage type, usage amount, usage start and end dates, and the unblended cost.

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage.  The `aws_cost_by_service_usage_type_daily` table provides a simplified view of cost for services in your account (or all linked accounts when run against the organization master), summarized by day, for the last year.  

Note that [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request will incur a cost of $0.01.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cost_by_service_usage_type_daily` table, you can use the `.inspect aws_cost_by_service_usage_type_daily` command in Steampipe.

**Key columns**:

- `service_name`: This column is essential as it specifies the AWS service name, allowing users to identify and analyze the cost associated with each service.
- `usage_type`: This column provides information about the type of usage for the AWS service. It is useful for understanding the specific service features contributing to the cost.
- `usage_start_date`: This column is important because it provides the start date for the usage period, enabling users to track and manage daily AWS service usage and costs.

## Examples

### Basic info

```sql
select
  service,
  usage_type,
  period_start,
  blended_cost_amount::numeric::money,
  unblended_cost_amount::numeric::money,
  amortized_cost_amount::numeric::money,
  net_unblended_cost_amount::numeric::money,
  net_amortized_cost_amount::numeric::money
from 
  aws_cost_by_service_usage_type_daily
order by
  service,
  period_start;
```



### Min, Max, and average daily unblended_cost_amount by service and usage type

```sql
select
  service,
  usage_type,
  min(unblended_cost_amount)::numeric::money as min,
  max(unblended_cost_amount)::numeric::money as max,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_service_usage_type_daily
group by
  service,
  usage_type
order by
  service,
  usage_type;
```

### Top 10 most expensive service usage type (by average daily unblended_cost_amount)

```sql
select
  service,
  usage_type,
  sum(unblended_cost_amount)::numeric::money as sum,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_service_usage_type_daily
group by
  service,
  usage_type
order by
  average desc
limit 10;
```


### Top 10 most expensive service usage type (by total daily unblended_cost_amount)

```sql
select
  service,
  usage_type,
  sum(unblended_cost_amount)::numeric::money as sum,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_service_usage_type_daily
group by
  service,
  usage_type
order by
  sum desc
limit 10;
```
