# Table: aws_cost_forecast_monthly

Amazon Cost Explorer helps you visualize, understand, and manage your AWS costs and usage.  The `aws_cost_forecast_monthly` table retrieves a forecast for how much Amazon Web Services predicts that you will spend each month over the next 12 months, based on your past costs.



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