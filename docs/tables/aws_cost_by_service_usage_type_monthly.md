---
title: "Steampipe Table: aws_cost_by_service_usage_type_monthly - Query AWS Cost Explorer Service using SQL"
description: "Allows users to query AWS Cost Explorer Service to get detailed cost data per service and usage type on a monthly basis."
folder: "Cost Explorer"
---

# Table: aws_cost_by_service_usage_type_monthly - Query AWS Cost Explorer Service using SQL

The AWS Cost Explorer Service is a tool that enables you to view and analyze your costs and usage. You can explore your AWS costs using an interface that lets you observe both your costs and usage patterns. It includes features that allow you to dive deeper into your cost and usage data to identify trends, pinpoint cost drivers, and detect anomalies.

## Table Usage Guide

The `aws_cost_by_service_usage_type_monthly` table in Steampipe provides you with information about the monthly cost data per service and usage type within AWS Cost Explorer Service. This table allows you, as a financial analyst or cloud cost manager, to query detailed cost data, including the service name, usage type, cost, and currency. You can utilize this table to gather insights on monthly AWS costs, such as cost per service, cost per usage type, and the total monthly cost. The schema outlines the various attributes of the cost data, including the service name, usage type, cost, and the currency used.

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage. The `aws_cost_by_service_usage_type_monthly` table provides you with a simplified view of cost for services in your account (or all linked accounts when run against the organization master), summarized by month, for the last year.  

**Important Notes**
- The [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request you make will incur a cost of $0.01.

## Examples

### Basic info
This query provides a comprehensive overview of your AWS service usage, allowing you to understand your monthly costs. By analyzing the cost and usage patterns, you can identify areas for potential cost savings and optimize your AWS utilization.

```sql+postgres
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

```sql+sqlite
select
  service,
  usage_type,
  period_start,
  cast(blended_cost_amount as decimal),
  cast(unblended_cost_amount as decimal),
  cast(amortized_cost_amount as decimal),
  cast(net_unblended_cost_amount as decimal),
  cast(net_amortized_cost_amount as decimal)
from 
  aws_cost_by_service_usage_type_monthly
order by
  service,
  period_start;
```

### Min, Max, and average monthly unblended_cost_amount by service and usage type
Gain insights into your AWS service usage by evaluating the minimum, maximum, and average monthly costs associated with each service and usage type. This helps in better understanding of your cloud spending patterns and can guide cost optimization efforts.

```sql+postgres
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

```sql+sqlite
select
  service,
  usage_type,
  min(cast(unblended_cost_amount as numeric)) as min,
  max(cast(unblended_cost_amount as numeric)) as max,
  avg(cast(unblended_cost_amount as numeric)) as average
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
Explore which services and usage types are the most costly on average per month, allowing for targeted cost reduction efforts. This analysis can help prioritize areas for cost optimization within your AWS services.

```sql+postgres
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

```sql+sqlite
select
  service,
  usage_type,
  sum(unblended_cost_amount) as sum,
  avg(unblended_cost_amount) as average
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
Discover the segments that are contributing the most to your monthly AWS costs. This query helps in identifying the top 10 service usage types that are incurring the highest costs, allowing you to better manage and optimize your resource usage.

```sql+postgres
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

```sql+sqlite
select
  service,
  usage_type,
  sum(unblended_cost_amount) as sum,
  avg(unblended_cost_amount) as average
from 
  aws_cost_by_service_usage_type_monthly
group by
  service,
  usage_type
order by
  sum desc
limit 10;
```