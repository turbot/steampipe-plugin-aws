---
title: "Steampipe Table: aws_cost_forecast_monthly - Query AWS Cost Explorer Cost Forecast using SQL"
description: "Allows users to query Cost Forecasts in AWS Cost Explorer for monthly cost predictions."
folder: "Cost Explorer"
---

# Table: aws_cost_forecast_monthly - Query AWS Cost Explorer Cost Forecast using SQL

The AWS Cost Explorer Cost Forecast is a feature of AWS that provides you with the ability to forecast your AWS costs. It uses your historical cost data to predict future expenses, enabling you to manage your budget more effectively. The forecasts are generated using machine learning algorithms and can be customized for different time periods, services, and tags.

## Table Usage Guide

The `aws_cost_forecast_monthly` table in Steampipe provides you with information about your monthly cost forecasts within AWS Cost Explorer. This table allows you, as a financial analyst or cloud cost manager, to query cost forecast details, including predicted costs, end and start dates, and associated metadata. You can utilize this table to gather insights on your future costs, such as predicted expenses for the next month, verification of cost trends, and more. The schema outlines the various attributes of your cost forecast, including the time period, value, and forecast results by time.

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage. The `aws_cost_forecast_monthly` table retrieves a forecast for how much Amazon Web Services predicts that you will spend each month over the next 12 months, based on your past costs.

**Important Notes**

- The [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request you make will incur a cost of $0.01.

## Examples

### Basic info
Assess the elements within your AWS cost forecast on a monthly basis to better understand your spending trends and budget accordingly. This query allows you to analyze your cost data over time, helping you to identify potential cost-saving opportunities and manage your AWS resources more effectively.

```sql+postgres

select 
   period_start,
   period_end,
   mean_value::numeric::money  
from 
  aws_cost_forecast_monthly
order by
  period_start;
```

```sql+sqlite
select 
   period_start,
   period_end,
   cast(mean_value as real) as mean_value
from 
  aws_cost_forecast_monthly
order by
  period_start;
```




###  Month on month forecasted growth
Gain insights into the monthly growth forecast by comparing the current month's mean value with the previous month's. This allows for a clear understanding of the growth percentage change, which can aid in future planning and budgeting.

```sql+postgres
with cost_data as (
  select
    period_start,
    mean_value as this_month,
    lag(mean_value,-1) over(order by period_start desc) as previous_month
  from 
    aws_cost_forecast_monthly
)
select
    period_start,
    this_month::numeric::money,
    previous_month::numeric::money,
    case 
      when previous_month = 0 and this_month = 0  then 0
      when previous_month = 0 then 999
      else round((100 * ( (this_month - previous_month) / previous_month))::numeric, 2) 
    end as percent_change
from
  cost_data
order by
  period_start;
```

```sql+sqlite
with cost_data as (
  select
    period_start,
    mean_value as this_month,
    lag(mean_value,-1) over(order by period_start desc) as previous_month
  from 
    aws_cost_forecast_monthly
)
select
    period_start,
    this_month,
    previous_month,
    case 
      when previous_month = 0 and this_month = 0  then 0
      when previous_month = 0 then 999
      else round((100 * ( (this_month - previous_month) / previous_month)), 2) 
    end as percent_change
from
  cost_data
order by
  period_start;
```