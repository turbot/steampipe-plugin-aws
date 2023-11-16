---
title: "Table: aws_cost_by_service_usage_type_monthly - Query AWS Cost Explorer Service using SQL"
description: "Allows users to query AWS Cost Explorer Service to get detailed cost data per service and usage type on a monthly basis."
---

# Table: aws_cost_by_service_usage_type_monthly - Query AWS Cost Explorer Service using SQL

The `aws_cost_by_service_usage_type_monthly` table in Steampipe provides information about the monthly cost data per service and usage type within AWS Cost Explorer Service. This table allows financial analysts and cloud cost managers to query detailed cost data, including the service name, usage type, cost, and currency. Users can utilize this table to gather insights on monthly AWS costs, such as cost per service, cost per usage type, and the total monthly cost. The schema outlines the various attributes of the cost data, including the service name, usage type, cost, and the currency used.

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage.  The `aws_cost_by_service_usage_type_monthly` table provides a simplified view of cost for services in your account (or all linked accounts when run against the organization master), summarized by month, for the last year.  

Note that [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request will incur a cost of $0.01.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cost_by_service_usage_type_monthly` table, you can use the `.inspect aws_cost_by_service_usage_type_monthly` command in Steampipe.

**Key columns**:

- `service_name`: The name of the AWS service. This column is important as it allows users to filter cost data by specific AWS services.
- `usage_type`: Describes the specific AWS operation for the usage. This column is useful for understanding the cost distribution across different types of usage within a service.
- `currency`: The currency in which the cost data is represented. This column is important for cost analysis across different currencies.

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
  aws_cost_by_service_usage_type_monthly
order by
  service,
  period_start;
```



### Min, Max, and average monthly unblended_cost_amount by service and usage type

```sql
select
  service,
  usage_type,
  min(unblended_cost_amount)::numeric::money as min,
  max(unblended_cost_amount)::numeric::money as max,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_service_usage_type_monthly
group by
  service,
  usage_type
order by
  service,
  usage_type;
```

### Top 10 most expensive service usage type (by average monthly unblended_cost_amount)

```sql
select
  service,
  usage_type,
  sum(unblended_cost_amount)::numeric::money as sum,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_service_usage_type_monthly
group by
  service,
  usage_type
order by
  average desc
limit 10;
```


### Top 10 most expensive service usage type (by total monthly unblended_cost_amount)

```sql
select
  service,
  usage_type,
  sum(unblended_cost_amount)::numeric::money as sum,
  avg(unblended_cost_amount)::numeric::money as average
from 
  aws_cost_by_service_usage_type_monthly
group by
  service,
  usage_type
order by
  sum desc
limit 10;
```
