---
title: "Table: aws_cost_forecast_daily - Query AWS Cost Explorer Daily Cost Forecast using SQL"
description: "Allows users to query AWS Cost Explorer's daily cost forecast data, providing insights into projected daily costs based on historical data."
---

# Table: aws_cost_forecast_daily - Query AWS Cost Explorer Daily Cost Forecast using SQL

The `aws_cost_forecast_daily` table in Steampipe provides information about daily cost forecasts within AWS Cost Explorer. This table allows financial analysts, DevOps engineers, and cloud administrators to query forecasted daily costs based on historical data. Users can utilize this table to gather insights on future costs, such as projected increases or decreases in expenses, cost trends, and more. The schema outlines the various attributes of the daily cost forecast, including the date, forecasted amount, and forecasted unit.

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage.  The `aws_cost_forecast_daily` table retrieves a forecast for how much Amazon Web Services predicts that you will spend each day over the next 4 months, based on your past costs.

Note that [pricing for the Cost Explorer API](https://aws.amazon.com/aws-cost-management/pricing/) is per API request - Each request will incur a cost of $0.01.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cost_forecast_daily` table, you can use the `.inspect aws_cost_forecast_daily` command in Steampipe.

### Key columns:

- `date`: This is the date for the forecasted cost. It is crucial for tracking cost trends over time and can be used to join this table with others containing date-specific information.
- `amount`: This column represents the forecasted cost amount for the specified date. It is essential for cost analysis and budgeting purposes.
- `unit`: This column indicates the unit of the forecasted cost amount (e.g., USD). It is necessary for understanding the currency in which costs are being forecasted.

## Examples

### Basic info

```sql

select 
   period_start,
   period_end,
   mean_value::numeric::money   
from 
  aws_cost_forecast_daily
order by
  period_start;
```

