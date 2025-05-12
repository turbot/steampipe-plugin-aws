---
title: "Steampipe Table: aws_cost_by_service_usage_type_daily - Query AWS Cost Explorer Service usage type daily using SQL"
description: "Allows users to query AWS Cost Explorer Service daily usage type to fetch detailed data about AWS service usage and costs."
folder: "Cost Explorer"
---

# Table: aws_cost_by_service_usage_type_daily - Query AWS Cost Explorer Service usage type daily using SQL

The AWS Cost Explorer Service usage type daily is a feature of AWS Cost Management that provides detailed information about your AWS costs, allowing you to visualize, understand, and manage your AWS costs and usage over time. This service provides data about your cost and usage in both tabular and graphical formats, with the ability to customize views and organize data to reflect your needs. The daily usage type specifically provides a granular view of costs incurred daily for each AWS service used.

## Table Usage Guide

The `aws_cost_by_service_usage_type_daily` table in Steampipe provides you with information about daily usage type and costs for each AWS service within AWS Cost Explorer. This table allows you, as a DevOps engineer, financial analyst, or cloud architect, to query daily-specific details, including usage amount, usage unit, and the corresponding service cost. You can utilize this table to gather insights on daily usage and costs, such as identifying high-cost services, tracking usage patterns, and managing your AWS expenses. The schema outlines the various attributes of the AWS service cost, including the service name, usage type, usage amount, usage start and end dates, and the unblended cost.

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage.  The `aws_cost_by_service_usage_type_daily` table provides you with a simplified view of cost for services in your account (or all linked accounts when run against the organization master), summarized by day, for the last year.  

**Important Notes**
- The [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request you make will incur a cost of $0.01.

## Examples

### Basic info
Explore your daily AWS service usage and costs, sorted by service and the start of the period. This can help you understand and manage your AWS expenses more effectively.

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
  aws_cost_by_service_usage_type_daily
order by
  service,
  period_start;
```

```sql+sqlite
select
  service,
  usage_type,
  period_start,
  CAST(blended_cost_amount AS NUMERIC) AS blended_cost_amount,
  CAST(unblended_cost_amount AS NUMERIC) AS unblended_cost_amount,
  CAST(amortized_cost_amount AS NUMERIC) AS amortized_cost_amount,
  CAST(net_unblended_cost_amount AS NUMERIC) AS net_unblended_cost_amount,
  CAST(net_amortized_cost_amount AS NUMERIC) AS net_amortized_cost_amount
from 
  aws_cost_by_service_usage_type_daily
order by
  service,
  period_start;
```



### Min, Max, and average daily unblended_cost_amount by service and usage type
Analyze your daily AWS service usage to understand the minimum, maximum, and average costs associated with each type of usage. This allows for more effective budget management and identification of potential cost-saving opportunities.

```sql+postgres
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

```sql+sqlite
select
  service,
  usage_type,
  min(unblended_cost_amount) as min,
  max(unblended_cost_amount) as max,
  avg(unblended_cost_amount) as average
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
Discover the segments that incur the highest average daily costs in your AWS services. This can help you identify areas where budget adjustments or cost optimizations might be necessary.

```sql+postgres
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

```sql+sqlite
select
  service,
  usage_type,
  sum(unblended_cost_amount) as sum,
  avg(unblended_cost_amount) as average
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
This query is used to analyze the most costly services in terms of daily usage. It helps in budget management by highlighting areas where costs are significantly high, thus aiding in cost optimization strategies.

```sql+postgres
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

```sql+sqlite
select
  service,
  usage_type,
  sum(unblended_cost_amount) as sum,
  avg(unblended_cost_amount) as average
from 
  aws_cost_by_service_usage_type_daily
group by
  service,
  usage_type
order by
  sum desc
limit 10;
```