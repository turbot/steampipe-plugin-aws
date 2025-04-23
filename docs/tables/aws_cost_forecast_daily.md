---
title: "Steampipe Table: aws_cost_forecast_daily - Query AWS Cost Explorer Daily Cost Forecast using SQL"
description: "Allows users to query AWS Cost Explorer's daily cost forecast data, providing insights into projected daily costs based on historical data."
folder: "Cost Explorer"
---

# Table: aws_cost_forecast_daily - Query AWS Cost Explorer Daily Cost Forecast using SQL

The AWS Cost Explorer Daily Cost Forecast is a feature of AWS that allows you to predict your future AWS costs based on your past spending. It uses machine learning algorithms to create a model of your past behavior and estimate your future costs. It provides an SQL interface for querying these forecasts, making it easy to integrate into your existing data analysis workflows.

## Table Usage Guide

The `aws_cost_forecast_daily` table in Steampipe provides you with daily cost forecasts within AWS Cost Explorer. This table allows you, as a financial analyst, DevOps engineer, or cloud administrator, to query forecasted daily costs based on historical data. You can utilize this table to gather insights on your future costs, such as projected increases or decreases in expenses, cost trends, and more. The schema outlines the various attributes of your daily cost forecast, including the date, forecasted amount, and forecasted unit.

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage. The `aws_cost_forecast_daily` table retrieves a forecast for how much Amazon Web Services predicts that you will spend each day over the next 4 months, based on your past costs.

**Important Notes**
- The [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request you make will incur a cost of $0.01.

## Examples

### Basic info
Explore the daily cost forecast for AWS, allowing you to understand and predict your expenditure over time. This can assist in budget planning and identifying potential cost-saving opportunities.

```sql+postgres

select 
   period_start,
   period_end,
   mean_value::numeric::money   
from 
  aws_cost_forecast_daily
order by
  period_start;
```

```sql+sqlite
select 
   period_start,
   period_end,
   cast(mean_value as decimal) as mean_value   
from 
  aws_cost_forecast_daily
order by
  period_start;
```