---
title: "Table: aws_cost_forecast_monthly - Query AWS Cost Explorer Cost Forecast using SQL"
description: "Allows users to query Cost Forecasts in AWS Cost Explorer for monthly cost predictions."
---

# Table: aws_cost_forecast_monthly - Query AWS Cost Explorer Cost Forecast using SQL

The `aws_cost_forecast_monthly` table in Steampipe provides information about the monthly cost forecasts within AWS Cost Explorer. This table allows financial analysts and cloud cost managers to query cost forecast details, including predicted costs, end and start dates, and associated metadata. Users can utilize this table to gather insights on future costs, such as predicted expenses for the next month, verification of cost trends, and more. The schema outlines the various attributes of the cost forecast, including the time period, value, and forecast results by time.

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage.  The `aws_cost_forecast_monthly` table retrieves a forecast for how much Amazon Web Services predicts that you will spend each month over the next 12 months, based on your past costs.

Note that [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request will incur a cost of $0.01.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cost_forecast_monthly` table, you can use the `.inspect aws_cost_forecast_monthly` command in Steampipe.

Key columns:

- `time_period_end`: This column provides the end date for the forecast period. It can be used to join with other tables that contain date-specific information.
- `time_period_start`: This column provides the start date for the forecast period. It can be used to join with other tables that contain date-specific information.
- `value_amount`: This column contains the predicted cost for the forecast period. It can be used to join with other tables that contain cost-related information.

## Examples

### Basic info

```sql

select 
   period_start,
   period_end,
   mean_value::numeric::money  
from 
  aws_cost_forecast_monthly
order by
  period_start;
```




###  Month on month forecasted growth

```sql
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